package handlers

import (
	"di/model"
	"di/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPipelines(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.UserReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user := &model.User{
			Username: req.Username,
			Password: req.Password,
		}

		err := services.UserService.Signin(user)

		if err != nil {
			log.Printf("Failed to log in user: %v\n", err.Error())
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, err := services.TokenService.NewFirstPairFromUser(context.Request.Context(), user)

		if err != nil {
			log.Printf("Failed to create tokens for user: %v\n", err.Error())

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})

			return
		}

		context.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func CreatePipeline(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.PipelineReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user, err := services.UserService.GetByUsername(req.User)

		if err != nil {
			log.Printf("Failed to log in user: %v\n", err.Error())
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// err := errors.NewNotFound("user", strconv.FormatUint(uint64(userId), 10))

		// context.JSON(err.Status(), gin.H{
		// 	"error": err,
		// })

		err := services.PipelineService.Create()

		if err != nil {
			log.Printf("Failed to log in user: %v\n", err.Error())
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, err := services.TokenService.NewFirstPairFromUser(context.Request.Context(), user)

		if err != nil {
			log.Printf("Failed to create tokens for user: %v\n", err.Error())

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})

			return
		}

		context.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
