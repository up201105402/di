package handlers

import (
	"di/model"
	"di/service"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

		for i := 0; i < len(pipelines); i++ {
			lastRun := getLastRun(services, pipelines[i].ID)
			pipelines[i].LastRun = lastRun
		}

		context.JSON(http.StatusOK, gin.H{
			"pipelines": pipelines,
		})
	}
}

func GetPipeline(services *service.Services) gin.HandlerFunc {
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

		pipeline.LastRun = getLastRun(services, pipeline.ID)

		context.JSON(http.StatusOK, gin.H{
			"pipeline": pipeline,
		})
	}
}

func UpsertPipeline(services *service.Services) gin.HandlerFunc {
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

		if req.ID != 0 {
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
		} else {
			serviceError := services.PipelineService.Create(user.ID, req.Name, req.Definition)

			if serviceError != nil {
				log.Printf("Failed to create in user: %v\n", err.Error())
				errorMessage := fmt.Sprint("Failed to create pipeline for user: %v\n", err.Error())
				err := errors.NewInternal()
				context.JSON(err.Status(), gin.H{
					"error": errorMessage,
				})
				return
			}
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

func getLastRun(services *service.Services, pipelineID uint) time.Time {
	runs, _ := services.RunService.GetByPipeline(pipelineID)

	var lastRun time.Time
	if runs != nil && len(runs) > 0 {
		lastRun = runs[0].LastRun

		for _, run := range runs {
			if run.LastRun.After(lastRun) {
				lastRun = run.LastRun
			}
		}
	}

	return lastRun
}
