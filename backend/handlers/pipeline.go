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
			errorMessage := fmt.Sprint("Failed to get pipelines for user with id %v: %v\n", user.ID, getError)
			log.Printf(errorMessage)
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
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		pipeline, getError := services.PipelineService.Get(uint(id))

		if getError != nil {
			errorMessage := fmt.Sprint("Failed to get pipeline with id %v: %v\n", pipelineId, getError)
			log.Printf(errorMessage)
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

func GetPipelineSchedule(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineId := context.Param("id")

		id, parseError := strconv.ParseUint(pipelineId, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		schedules, getError := services.PipelineService.GetPipelineSchedules(uint(id))

		if getError != nil {
			errorMessage := fmt.Sprint("Failed to get pipeline's schedules with id %s: %v\n", pipelineId, getError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"schedules": schedules,
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
				errorMessage := fmt.Sprint("Failed to update pipeline with id &d: %v\n", req.ID, err.Error())
				log.Printf(errorMessage)
				err := errors.NewInternal()
				context.JSON(err.Status(), gin.H{
					"error": errorMessage,
				})
				return
			}
		} else {
			serviceError := services.PipelineService.Create(user.ID, req.Name, req.Definition)

			if serviceError != nil {
				errorMessage := fmt.Sprint("Failed to create pipeline for user: %v\n", err.Error())
				log.Print(errorMessage)
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

func CreatePipelineSchedule(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineID := context.Param("id")

		id, parseError := strconv.ParseUint(pipelineID, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		var req model.PipelineScheduleReq

		if ok := bindData(context, &req); !ok {
			return
		}

		if req.CronExpression != "" || req.UniqueOcurrence.Year() > 1 {
			user, err := getUser(context)
			if err != nil {
				context.JSON(err.Status(), gin.H{
					"error": err.Error(),
				})
			}

			pipeline, pipelineErr := services.PipelineService.Get(uint(id))

			if pipelineErr != nil {
				log.Printf("Failed to get pipeline with id %d\n", req.ID)
				err := errors.NewNotFound("pipeline", string(req.ID))
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			if user.ID != pipeline.UserID {
				msg := fmt.Sprintf("Pipeline %s is not owned by user %s\n", pipeline.Name, user.Username)
				log.Printf(msg)
				err := errors.NewAuthorization(msg)
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			createError := services.PipelineService.CreatePipelineSchedule(pipeline.ID, req.UniqueOcurrence, req.CronExpression)
			if createError != nil {
				errorMessage := fmt.Sprint("Failed to get pipeline's schedules with id %s: %v\n", pipelineID, createError)
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

func UploadPipelineFile(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		// Upload the file to specific dst.
		fileUploadDir, exists := os.LookupEnv("FILE_UPLOAD_DIR")

		if !exists {
			errorMessage := fmt.Sprint("Upload directory is not defined")
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
		}

		pipelineID := context.Param("id")

		_, parseError := strconv.ParseUint(pipelineID, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		file, _ := context.FormFile("file")
		file.Filename = "file_" + fmt.Sprintf("%d", time.Now().Unix())
		log.Println(file.Filename)

		err := context.SaveUploadedFile(file, fileUploadDir+"pipelines/"+pipelineID+"/"+file.Filename)

		if err != nil {
			errorMessage := fmt.Sprintf("Failed to save uploaded file: %v", err.Error())
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
		}

		context.JSON(http.StatusOK, gin.H{
			"filename": file.Filename,
		})
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
			errorMessage := fmt.Sprint("Failed to get pipelines for user with id %v: %v\n", user.ID, getError)
			log.Printf(errorMessage)
			err := errors.NewNotFound("pipeline", string(req.ID))
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		if pipeline.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.PipelineService.Delete(req.ID)

		if deleteError != nil {
			errorMessage := fmt.Sprint("Failed to delete pipeline with id %v: %v\n", req.ID, err.Error())
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

func DeletePipelineSchedule(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineId := context.Param("id")

		pipelineID, parseError := strconv.ParseUint(pipelineId, 10, 64)

		if parseError != nil {
			errorMessage := fmt.Sprint("Failed to convert pipelineId into uint: %v\n", parseError)
			log.Printf(errorMessage)
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		var req model.PipelineScheduleReq

		if ok := bindData(context, &req); !ok {
			return
		}

		pipeline, getError := services.PipelineService.Get(uint(pipelineID))

		if getError != nil {
			errorMessage := fmt.Sprint("Failed to get pipeline with id %v\n", getError)
			log.Printf(errorMessage)
			err := errors.NewNotFound("pipeline", string(pipelineID))
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

		if pipeline.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", pipeline.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.PipelineService.DeletePipelineSchedule(req.ID)

		if deleteError != nil {
			errorMessage := fmt.Sprint("Failed to delete pipeline schedule with id %v: %v\n", req.ID, err.Error())
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
