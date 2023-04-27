package handlers

import (
	"di/model"
	"di/service"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPipelines(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		pipelines, getError := services.PipelineService.GetByOwner(user.ID)

		if getError != nil {
			log.Printf("Failed to get pipelines for user with id %v: %v\n", user.ID, getError)
			errorMessage := fmt.Sprint("Failed to get pipelines for user with id %v: %v\n", user.ID, getError)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"pipelines": pipelines,
		})
	}
}

func CreatePipeline(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.PipelineReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		createError := services.PipelineService.Create(user.ID, req.Name, req.Definition)

		if createError != nil {
			log.Printf("Failed to create in user: %v\n", err.Error())
			errorMessage := fmt.Sprint("Failed to create pipeline for user: %v\n", err.Error())
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func DeletePipeline(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.PipelineReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		pipeline, getError := services.PipelineService.Get(req.ID)

		if getError != nil {
			log.Printf("Failed to get pipeline with id %v\n", getError)
			errorMessage := fmt.Sprint("Failed to get pipelines for user with id %v: %v\n", user.ID, getError)
			err := errors.NewNotFound("pipeline", string(req.ID))
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		if pipeline.UserID != user.ID {
			log.Printf("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.PipelineService.Delete(req.ID)

		if deleteError != nil {
			log.Printf("Failed to delete pipeline with id %v: %v\n", req.ID, err.Error())
			errorMessage := fmt.Sprint("Failed to delete pipeline with id %v: %v\n", req.ID, err.Error())
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func getUser(context *gin.Context) (*model.User, *errors.Error) {
	contextUser, exists := context.Get("user")

	if !exists {
		err := errors.NewNotFound("user", "")
		return nil, err
	}

	return contextUser.(*model.User), nil
}