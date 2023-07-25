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
	trainedService := service.NewTrainedService(dbConnection, client, i18n)
	runService := service.NewRunService(dbConnection, client, i18n, &pipelineService, &stepTypeService, &trainedService)
	taskService := service.NewTaskService(i18n, &stepTypeService, &runService)
	datasetService := service.NewDatasetService(dbConnection, client, i18n)
	trainerService := service.NewTrainerService(dbConnection, client, i18n)
	testerService := service.NewTesterService(dbConnection, client, i18n)

	services := &service.Services{
		UserService:     service.NewUserService(dbConnection, i18n),
		PipelineService: pipelineService,
		RunService:      runService,
		TokenService:    service.NewTokenService(tokenServiceConfig, i18n),
		DatasetService:  datasetService,
		TrainerService:  trainerService,
		TesterService:   testerService,
		TrainedService:  trainedService,
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

	logsDir, exists := os.LookupEnv("RUN_LOGS_DIR")

	if !exists {
		panic("RUN_LOGS_DIR is not defined")
	}

	filesDir, exists := os.LookupEnv("FILE_UPLOAD_DIR")

	if !exists {
		panic("FILE_UPLOAD_DIR is not defined")
	}

	router.StaticFS("/work", http.Dir(workDir))
	router.StaticFS("/logs", http.Dir(logsDir))
	router.StaticFS("/files", http.Dir(filesDir))

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
	pipelineAPI.POST("/", middleware.Auth(services.TokenService, I18n), handlers.UpsertPipeline(services))
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
	runAPI.POST("/resume/:runID", middleware.Auth(services.TokenService, I18n), handlers.ResumeRun(services, I18n))

	runResultsAPI := router.Group("/api/runresults")
	runResultsAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.FindRunResulstById(services, I18n))
	runResultsAPI.GET("/:id/log", middleware.Auth(services.TokenService, I18n), handlers.GetLogTail(services, I18n))

	feedbackAPI := router.Group("/api/feedback")
	feedbackAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.FindRunFeedbackQueriesByRunId(services, I18n))
	feedbackAPI.GET("/:id/query/:queryId", middleware.Auth(services.TokenService, I18n), handlers.FindRunFeedbackQueryById(services, I18n))
	feedbackAPI.POST("/:id", middleware.Auth(services.TokenService, I18n), handlers.SubmitRunFeedback(services, I18n))

	databasetAPI := router.Group("/api/dataset")
	databasetAPI.GET("", middleware.Auth(services.TokenService, I18n), handlers.GetDatasets(services))
	databasetAPI.POST("", middleware.Auth(services.TokenService, I18n), handlers.CreateDataset(services))
	databasetAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.GetDataset(services, I18n))
	databasetAPI.POST("/:id/file", middleware.Auth(services.TokenService, I18n), handlers.UploadDatasetScript(services, I18n))
	databasetAPI.DELETE("", middleware.Auth(services.TokenService, I18n), handlers.DeleteDataset(services, I18n))

	trainerAPI := router.Group("/api/trainer")
	trainerAPI.GET("", middleware.Auth(services.TokenService, I18n), handlers.GetTrainers(services))
	trainerAPI.POST("", middleware.Auth(services.TokenService, I18n), handlers.CreateTrainer(services))
	trainerAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.GetTrainer(services, I18n))
	trainerAPI.POST("/:id/file", middleware.Auth(services.TokenService, I18n), handlers.UploadTrainerScript(services, I18n))
	trainerAPI.DELETE("", middleware.Auth(services.TokenService, I18n), handlers.DeleteTrainer(services, I18n))

	testerAPI := router.Group("/api/tester")
	testerAPI.GET("", middleware.Auth(services.TokenService, I18n), handlers.GetTesters(services))
	testerAPI.POST("", middleware.Auth(services.TokenService, I18n), handlers.CreateTester(services))
	testerAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.GetTester(services, I18n))
	testerAPI.POST("/:id/file", middleware.Auth(services.TokenService, I18n), handlers.UploadTesterScript(services, I18n))
	testerAPI.DELETE("", middleware.Auth(services.TokenService, I18n), handlers.DeleteTester(services, I18n))

	trainedAPI := router.Group("/api/trained")
	trainedAPI.GET("", middleware.Auth(services.TokenService, I18n), handlers.GetTrainedModels(services))
	trainedAPI.POST("", middleware.Auth(services.TokenService, I18n), handlers.CreateTrained(services))
	trainedAPI.GET("/:id", middleware.Auth(services.TokenService, I18n), handlers.GetTrained(services, I18n))
	trainedAPI.POST("/:id/file", middleware.Auth(services.TokenService, I18n), handlers.UploadTrainedScript(services, I18n))
	trainedAPI.DELETE("", middleware.Auth(services.TokenService, I18n), handlers.DeleteTrained(services, I18n))

	return router
}
