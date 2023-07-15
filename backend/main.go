package main

import (
	"di/handlers"
	"di/middleware"
	"di/service"
	"di/util"
	"net/http"
	"os"

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

	client := util.GetAsynqClient()
	defer client.Close()

	pipelineService := service.NewPipelineService(dbConnection, client)
	pipelineService.SyncAsyncTasks()
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

	port, exists := os.LookupEnv("WEB_SERVER_PORT")

	if !exists {
		port = "8001"
	}

	r.Run(":" + port)
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
	pipelineAPI.GET("/:id/schedule", middleware.Auth(services.TokenService), handlers.GetPipelineSchedule(services))
	pipelineAPI.POST("/:id/schedule", middleware.Auth(services.TokenService), handlers.CreatePipelineSchedule(services))
	pipelineAPI.POST("/:id/file", middleware.Auth(services.TokenService), handlers.UploadPipelineFile(services))
	pipelineAPI.POST("/:id", middleware.Auth(services.TokenService), handlers.UpsertPipeline(services))
	pipelineAPI.DELETE("", middleware.Auth(services.TokenService), handlers.DeletePipeline(services))
	pipelineAPI.DELETE("/:id/schedule", middleware.Auth(services.TokenService), handlers.DeletePipelineSchedule(services))

	runAPI := router.Group("/api/run")
	runAPI.GET("", middleware.Auth(services.TokenService), handlers.GetRuns(services))
	runAPI.GET("/:id", middleware.Auth(services.TokenService), handlers.FindRunsByPipeline(services))
	runAPI.POST("/:id", middleware.Auth(services.TokenService), handlers.CreateRun(services))
	runAPI.POST("/execute/:runID", middleware.Auth(services.TokenService), handlers.ExecuteRun(services))

	runResultsAPI := router.Group("/api/runresults")
	runResultsAPI.GET("/:id", middleware.Auth(services.TokenService), handlers.FindRunResulstById(services))

	return router
}
