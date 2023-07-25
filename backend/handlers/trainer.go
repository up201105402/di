package handlers

import (
	"di/model"
	"di/service"
	"di/util"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetTrainers(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		trainers, getError := services.TrainerService.GetByOwner(user.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		fileUploadDir := os.Getenv("FILE_UPLOAD_DIR")

		for index, trainer := range trainers {
			if trainer.Path != "" {
				path := "/files/" + strings.Split(trainer.Path, fileUploadDir)[1]
				trainers[index].Path = path
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"trainers": trainers,
		})
	}
}

func GetTrainer(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		trainerID := context.Param("id")

		id, parseError := strconv.ParseUint(trainerID, 10, 64)

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

		trainer, getError := services.TrainerService.Get(uint(id))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"model": trainer,
		})
	}
}

func CreateTrainer(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.TrainerReq

		if ok := util.BindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		if req.ID != 0 {
			trainer, err := services.TrainerService.Get(req.ID)

			if err != nil {
				err := errors.NewNotFound(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			trainer.Path = req.Path
			err = services.TrainerService.Update(trainer)

			if err != nil {
				log.Printf(err.Error())
				err := errors.NewInternal(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		} else {
			serviceError := services.TrainerService.Create(user.ID, req.Name)

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

func UploadTrainerScript(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
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

		trainerId := context.Param("id")
		trainerID, parseError := strconv.ParseUint(trainerId, 10, 64)

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

		trainer, getError := services.TrainerService.Get(uint(trainerID))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		trainerUploadDir := fileUploadDir + "trainers/" + trainerId + "/"
		if err := os.MkdirAll(trainerUploadDir, os.ModePerm); err != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "os.cmd.mkdir.dir.failed",
				TemplateData: map[string]interface{}{
					"Path":   trainerUploadDir,
					"Reason": err.Error(),
				},
				PluralCount: 1,
			})

			log.Println(errMessage)
			err := errors.NewInternal(errMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		file, _ := context.FormFile("file")
		log.Println(file.Filename)

		filePath := filepath.Join(trainerUploadDir, file.Filename)
		err := context.SaveUploadedFile(file, filePath)

		if err != nil {
			log.Printf(err.Error())
			err := errors.NewInternal(err.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
		}

		trainer.Path = filePath
		err = services.TrainerService.Update(trainer)

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

func UpdateTrainer(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.TrainerReq

		if ok := util.BindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		trainer, getError := services.TrainerService.Get(req.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if trainer.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.TrainerService.Delete(req.ID)

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

func DeleteTrainer(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.TrainerReq

		if ok := util.BindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		trainer, getError := services.TrainerService.Get(req.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if trainer.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.TrainerService.Delete(req.ID)

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
