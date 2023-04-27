package main

import (
	"di/handlers"
	"di/middleware"
	"di/model"
	"di/service"
	"di/util"
	"log"
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

	services := &service.Services{
		UserService:     service.NewUserService(dbConnection),
		PipelineService: service.NewPipelineService(dbConnection),
		TokenService:    service.NewTokenService(tokenServiceConfig),
	}

	r := setupRouter(services)

	r.Run(":8001")
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
	userAPI.POST("/signout", middleware.Auth(services.TokenService), handlers.SignOut(services))

	pipelineAPI := router.Group("/api/pipeline")
	pipelineAPI.GET("", middleware.Auth(services.TokenService), handlers.GetPipelines(services))
	pipelineAPI.POST("", middleware.Auth(services.TokenService), handlers.CreatePipeline(services))
	pipelineAPI.DELETE("", middleware.Auth(services.TokenService), handlers.DeletePipeline(services))

	return router
}
