package main

import (
	"di/handlers"
	"di/middleware"
	"di/service"
	"di/util"
	"net/http"

	"github.com/gin-gonic/gin"
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
	stepTypeService := service.NewNodeService()
	runService := service.NewRunService(dbConnection, client, &pipelineService, &stepTypeService)
	taskService := service.NewTaskService(&stepTypeService, &runService)

	services := &service.Services{
		UserService:     service.NewUserService(dbConnection),
		PipelineService: pipelineService,
		RunService:      runService,
		TokenService:    service.NewTokenService(tokenServiceConfig),
	}

	r := setupRouter(services)

	go taskService.SetupAsynqWorker()

	r.Run(":8001")
}

func setupRouter(services *service.Services) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "../client/public")
	router.Static("/data-sources", "../client/public/data-sources")
	router.Static("/assets", "../client/src/assets")

	router.LoadHTMLFiles("../client/index.html")

	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	userAPI := router.Group("/api/user")
	userAPI.POST("/login", handlers.LogIn(services))
	userAPI.POST("/signup", handlers.SignUp(services))
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
