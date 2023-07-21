package main

import (
	"di/handlers"
	"di/middleware"
	"di/service"
	"di/util"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

	i18n := util.GetI18nLocalizer()

	pipelineService := service.NewPipelineService(dbConnection, client, i18n)
	pipelineService.SyncAsyncTasks()
	stepTypeService := service.NewNodeService(i18n)
	runService := service.NewRunService(dbConnection, client, i18n, &pipelineService, &stepTypeService)
	taskService := service.NewTaskService(i18n, &stepTypeService, &runService)

	services := &service.Services{
		UserService:     service.NewUserService(dbConnection, i18n),
		PipelineService: pipelineService,
		RunService:      runService,
		TokenService:    service.NewTokenService(tokenServiceConfig, i18n),
	}

	r := setupRouter(services, i18n)

	go taskService.SetupAsynqWorker()

	port, exists := os.LookupEnv("WEB_SERVER_PORT")

	if !exists {
		port = "8001"
	}

	r.Run(":" + port)
}

func setupRouter(services *service.Services, I18n *i18n.Localizer) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "../client/public")
	router.Static("/assets", "../client/src/assets")

	workDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		panic("PIPELINES_WORK_DIR is not defined")
	}

	router.StaticFS("/work", http.Dir(workDir))

	router.LoadHTMLFiles("../client/index.html")

	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	userAPI := router.Group("/api/user")
	userAPI.POST("/", handlers.EditUser(services, I18n))
	userAPI.POST("/login", handlers.LogIn(services))
	userAPI.POST("/signup", handlers.SignUp(services))
	userAPI.POST("/tokens", handlers.NewAccessToken(services))
	userAPI.POST("/signout", middleware.Auth(services.TokenService, I18n), handlers.SignOut(services))

	pipelineAPI := router.Group("/api/pipeline")
	pipelineAPI.GET("", middleware.Auth(services.TokenService, I18n), handlers.GetPipelines(services))
	pipelineAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.GetPipeline(services, I18n))
	pipelineAPI.GET("/:id/schedule", middleware.Auth(services.TokenService, I18n), handlers.GetPipelineSchedule(services, I18n))
	pipelineAPI.POST("/:id/schedule", middleware.Auth(services.TokenService, I18n), handlers.CreatePipelineSchedule(services, I18n))
	pipelineAPI.POST("/:id/file", middleware.Auth(services.TokenService, I18n), handlers.UploadPipelineFile(services, I18n))
	pipelineAPI.POST("/:id", middleware.Auth(services.TokenService, I18n), handlers.UpsertPipeline(services))
	pipelineAPI.DELETE("", middleware.Auth(services.TokenService, I18n), handlers.DeletePipeline(services))
	pipelineAPI.DELETE("/:id/schedule", middleware.Auth(services.TokenService, I18n), handlers.DeletePipelineSchedule(services, I18n))

	runAPI := router.Group("/api/run")
	runAPI.GET("", middleware.Auth(services.TokenService, I18n), handlers.GetRuns(services))
	runAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.FindRunsByPipeline(services, I18n))
	runAPI.POST("/:id", middleware.Auth(services.TokenService, I18n), handlers.CreateRun(services, I18n))
	runAPI.POST("/execute/:runID", middleware.Auth(services.TokenService, I18n), handlers.ExecuteRun(services, I18n))

	runResultsAPI := router.Group("/api/runresults")
	runResultsAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.FindRunResulstById(services, I18n))

	feedbackAPI := router.Group("/api/feedback")
	feedbackAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.FindRunFeedbackQueriesById(services, I18n))

	return router
}
