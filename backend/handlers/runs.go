package handlers

import (
	"di/model"
	"di/service"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
			errorMessage := fmt.Sprint("Failed to get pipelines for user with id %v: %v\n", user.ID, getError)
			log.Printf(errorMessage)
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

func FindRunsByPipeline(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		pipelineID, parseError := strconv.ParseUint(id, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		pipeline, serviceError := services.PipelineService.Get(uint(pipelineID))

		if serviceError != nil {
			errorMessage := fmt.Sprint("Failed to get pipeline with id %v: %v\n", pipelineID, serviceError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		runs, getError := services.RunService.GetByPipeline(pipeline.ID)

		if serviceError != nil {
			errorMessage := fmt.Sprint("Failed to get runs for pipeline with id %v: %v\n", pipelineID, getError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"runs": runs,
		})
	}
}

func CreateRun(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		pipelineID, parseError := strconv.ParseUint(id, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		var req model.CreateRunReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		pipeline, serviceError := services.PipelineService.Get(uint(pipelineID))

		if serviceError != nil {
			err := errors.NewNotFound("pipeline", string(pipelineID))
			log.Printf(err.Message)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if pipeline.User.ID != user.ID {
			errorMessage := fmt.Sprint("Failed to create run for pipeline %d with user: %v\n", pipeline.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		_, serviceError = services.RunService.Create(*pipeline)

		if serviceError != nil {
			errorMessage := fmt.Sprint("Failed to create run for pipeline: %v\n", err.Error())
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		if req.Execute {
			serviceError := services.RunService.Execute(pipeline.ID)

			if serviceError != nil {
				errorMessage := fmt.Sprint("Failed to execute run for pipeline: %v\n", err.Error())
				log.Printf(errorMessage)
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

func ExecuteRun(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("runID")

		runID, parseError := strconv.ParseUint(id, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert runID into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		run, serviceError := services.RunService.Get(uint(runID))

		if serviceError != nil {
			err := errors.NewNotFound("pipeline", string(runID))
			log.Printf(err.Message)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if run.Pipeline.User.ID != user.ID {
			errorMessage := fmt.Sprint("Failed to execute pipeline %d with user: %v\n", run.Pipeline.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		serviceError = services.RunService.Execute(run.ID)

		if serviceError != nil {
			errorMessage := fmt.Sprint("Failed to execute run for pipeline: %v\n", err.Error())
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func FindRunResulstById(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		runID, parseError := strconv.ParseUint(id, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert runId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		run, serviceError := services.RunService.Get(uint(runID))

		if serviceError != nil {
			errorMessage := fmt.Sprint("Failed to get run with id %v: %v\n", runID, serviceError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		runStepStatuses, getError := services.RunService.FindRunStepStatusesByRun(run.ID)

		if serviceError != nil {
			errorMessage := fmt.Sprint("Failed to get runs for pipeline with id %v: %v\n", runID, getError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		runLogsDir, exists := os.LookupEnv("RUN_LOGS_DIR")

		if !exists {
			errorMessage := fmt.Sprint("Run logs directory is not defined!")
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		logFileName, exists := os.LookupEnv("RUN_LOG_FILE_NAME")

		if !exists {
			errorMessage := fmt.Sprint("Run log file name is not defined!")
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		runLogDir := runLogsDir + "/pipelines/" + fmt.Sprint(run.PipelineID) + "/" + fmt.Sprint(run.ID) + "/"
		logFile, _ := os.ReadFile(runLogDir + logFileName)

		context.JSON(http.StatusOK, gin.H{
			"run":             run,
			"runStepStatuses": runStepStatuses,
			"log":             string(logFile),
		})
	}
}
