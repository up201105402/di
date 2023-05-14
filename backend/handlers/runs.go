package handlers

import (
	"di/model"
	"di/service"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRuns(services *service.Services) gin.HandlerFunc {
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

func FindByPipeline(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineId := context.Param("id")

		id, parseError := strconv.ParseUint(pipelineId, 10, 64)

		if parseError != nil {
			log.Printf("Failed to convert pipelineId into uint: %v\n", parseError)
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		pipeline, getError := services.PipelineService.Get(uint(id))

		if getError != nil {
			log.Printf("Failed to get pipeline with id %v: %v\n", pipelineId, getError)
			errorMessage := fmt.Sprint("Failed to get pipeline with id %v: %v\n", pipelineId, getError)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"pipeline": pipeline,
		})
	}
}

func CreateRun(services *service.Services) gin.HandlerFunc {
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

		serviceError := services.PipelineService.Create(user.ID, req.Name, req.Definition)

		if req.ID != 0 && req.Definition != "" {
			pipeline, err := services.PipelineService.Get(req.ID)

			if err != nil {
				log.Printf("Failed to get pipeline with id %d\n", req.ID)
				err := errors.NewNotFound("pipeline", string(req.ID))
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			pipeline.Definition = req.Definition
			err = services.PipelineService.Update(pipeline)

			if err != nil {
				log.Printf("Failed to update pipeline with id %d\n", req.ID)
				errorMessage := fmt.Sprint("Failed to update pipeline with id &d: %v\n", req.ID, err.Error())
				err := errors.NewInternal()
				context.JSON(err.Status(), gin.H{
					"error": errorMessage,
				})
				return
			}
		}

		if serviceError != nil {
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
