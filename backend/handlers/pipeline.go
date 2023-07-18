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
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
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

func GetPipeline(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineId := context.Param("id")

		id, parseError := strconv.ParseUint(pipelineId, 10, 64)

		if parseError != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "sys.parsing.string.uint",
				TemplateData: map[string]interface{}{
					"Reason": parseError.Error(),
				},
				PluralCount: 1,
			})
			log.Printf(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		pipeline, getError := services.PipelineService.Get(uint(id))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		pipeline.LastRun = getLastRun(services, pipeline.ID)

		context.JSON(http.StatusOK, gin.H{
			"pipeline": pipeline,
		})
	}
}

func GetPipelineSchedule(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineId := context.Param("id")

		id, parseError := strconv.ParseUint(pipelineId, 10, 64)

		if parseError != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "sys.parsing.string.uint",
				TemplateData: map[string]interface{}{
					"Reason": parseError.Error(),
				},
				PluralCount: 1,
			})
			log.Printf(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		schedules, getError := services.PipelineService.GetPipelineSchedules(uint(id))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
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
				err := errors.NewNotFound(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			pipeline.Definition = req.Definition
			err = services.PipelineService.Update(pipeline)

			if err != nil {
				log.Printf(err.Error())
				err := errors.NewInternal(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		} else {
			serviceError := services.PipelineService.Create(user.ID, req.Name, req.Definition)

			if serviceError != nil {
				log.Print(serviceError.Error())
				err := errors.NewInternal(serviceError.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func CreatePipelineSchedule(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineID := context.Param("id")

		id, parseError := strconv.ParseUint(pipelineID, 10, 64)

		if parseError != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "sys.parsing.string.uint",
				TemplateData: map[string]interface{}{
					"Reason": parseError.Error(),
				},
				PluralCount: 1,
			})
			log.Printf(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
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
				err := errors.NewNotFound(pipelineErr.Error())
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
				log.Printf(createError.Error())
				err := errors.NewInternal(createError.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func UploadPipelineFile(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		// Upload the file to specific dst.
		fileUploadDir, exists := os.LookupEnv("FILE_UPLOAD_DIR")

		if !exists {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "env.variable.find.failed",
				TemplateData: map[string]interface{}{
					"Name": "FILE_UPLOAD_DIR",
				},
				PluralCount: 1,
			})

			log.Printf(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
		}

		pipelineID := context.Param("id")

		_, parseError := strconv.ParseUint(pipelineID, 10, 64)

		if parseError != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "sys.parsing.string.uint",
				TemplateData: map[string]interface{}{
					"Reason": parseError.Error(),
				},
				PluralCount: 1,
			})
			log.Printf(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		file, _ := context.FormFile("file")
		file.Filename = "file_" + fmt.Sprintf("%d", time.Now().Unix())
		log.Println(file.Filename)

		err := context.SaveUploadedFile(file, fileUploadDir+"pipelines/"+pipelineID+"/"+file.Filename)

		if err != nil {
			log.Printf(err.Error())
			err := errors.NewInternal(err.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
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
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
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
			log.Printf(deleteError.Error())
			err := errors.NewInternal(deleteError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func DeletePipelineSchedule(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		pipelineId := context.Param("id")

		pipelineID, parseError := strconv.ParseUint(pipelineId, 10, 64)

		if parseError != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "sys.parsing.string.uint",
				TemplateData: map[string]interface{}{
					"Reason": parseError.Error(),
				},
				PluralCount: 1,
			})
			log.Printf(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		var req model.PipelineScheduleReq

		if ok := bindData(context, &req); !ok {
			return
		}

		pipeline, getError := services.PipelineService.Get(uint(pipelineID))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
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
			log.Printf(deleteError.Error())
			err := errors.NewInternal(deleteError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{})
	}
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
