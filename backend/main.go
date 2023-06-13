package main

import (
	"di/handlers"
	"di/middleware"
	"di/model"
	"di/service"
	"di/tasks"
	"di/util"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func main() {
	dbConnection, err := util.ConnectToDB()

	if err != nil {
		panic("Failed to connect to the database")
	}

	redisClient, err := util.ConnectToRedis()

	if err != nil {
		panic("Failed to connect to Redis")
	}

	if err := util.CreateOrUpdateSchema(dbConnection); err != nil {
		panic("Failed to connect to Redis")
	}

	tokenServiceConfig, err := service.GetTokenServiceConfig(redisClient)

	if err != nil {
		panic("Failed to get Token Configuration")
	}

	client := util.ConnectToAsynq()
	defer client.Close()

	pipelineService := service.NewPipelineService(dbConnection)
	stepTypeService := service.NewStepService()

	services := &service.Services{
		UserService:     service.NewUserService(dbConnection),
		PipelineService: pipelineService,
		RunService:      service.NewRunService(dbConnection, client, &pipelineService, &stepTypeService),
		TokenService:    service.NewTokenService(tokenServiceConfig),
	}

	r := setupRouter(services)

	go setupAsynqWorker()

	if err != nil {
		panic("Failed to config Asynq")
	}

	r.Run(":8001")
}

func setupAsynqWorker() error {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisConnection := asynq.RedisClientOpt{
		Addr: redisHost + ":" + redisPort,
	}

	worker := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"runs": 1,
		},
	})

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		tasks.RunPipelineTask,
		tasks.HandleRunPipelineTask,
	)

	if err := worker.Run(mux); err != nil {
		return err
	}

	return nil
}

func setupRouter(services *service.Services) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "../client/public")
	router.Static("/data-sources", "../client/public/data-sources")
	router.Static("/assets", "../client/src/assets")

	router.LoadHTMLFiles("../client/index.html")

	router.GET("/", func(context *gin.Context) {
		username, exists := context.Get("user")

		if exists {
			userId := username.(*model.User).ID

			user, err := services.UserService.Get(userId)

			if err != nil {
				log.Printf("Unable to find user: %d\n%v", userId, err)
				// err := errors.NewNotFound("user", strconv.FormatUint(uint64(userId), 10))

				// context.JSON(err.Status(), gin.H{
				// 	"error": err,
				// })

				// return
			}

			context.HTML(http.StatusOK, "index.html", gin.H{
				"user": user,
			})

			return
		}

		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	userAPI := router.Group("/api/user")
	userAPI.POST("/login", handlers.LogIn(services))
	userAPI.POST("/signup", middleware.Auth(services.TokenService), handlers.SignUp(services))
	userAPI.POST("/tokens", handlers.NewAccessToken(services))
	userAPI.POST("/signout", middleware.Auth(services.TokenService), handlers.SignOut(services))

	pipelineAPI := router.Group("/api/pipeline")
	pipelineAPI.GET("", middleware.Auth(services.TokenService), handlers.GetPipelines(services))
	pipelineAPI.GET("/:id", middleware.Auth(services.TokenService), handlers.GetPipeline(services))
	pipelineAPI.POST("/:id", middleware.Auth(services.TokenService), handlers.UpsertPipeline(services))
	pipelineAPI.DELETE("", middleware.Auth(services.TokenService), handlers.DeletePipeline(services))

	runAPI := router.Group("/api/run")
	runAPI.GET("", middleware.Auth(services.TokenService), handlers.GetRuns(services))
	runAPI.GET("/:id", middleware.Auth(services.TokenService), handlers.FindByPipeline(services))
	runAPI.POST("/:id", middleware.Auth(services.TokenService), handlers.CreateRun(services))
	runAPI.POST("/execute/:runID", middleware.Auth(services.TokenService), handlers.ExecuteRun(services))

	return router
}
