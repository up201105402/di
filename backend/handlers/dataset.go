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

func GetDatasets(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		datasets, getError := services.DatasetService.GetByOwner(user.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		fileUploadDir := os.Getenv("FILE_UPLOAD_DIR")

		for index, dataset := range datasets {
			if dataset.Path != "" {
				path := "/files/" + strings.Split(dataset.Path, fileUploadDir)[1]
				datasets[index].Path = path
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"datasets": datasets,
		})
	}
}

func GetDataset(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		datasetID := context.Param("id")

		id, parseError := strconv.ParseUint(datasetID, 10, 64)

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

		dataset, getError := services.DatasetService.Get(uint(id))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"dataset": dataset,
		})
	}
}

func CreateDataset(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.DatasetReq

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
			dataset, err := services.DatasetService.Get(req.ID)

			if err != nil {
				err := errors.NewNotFound(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			dataset.Path = req.Path
			err = services.DatasetService.Update(dataset)

			if err != nil {
				log.Printf(err.Error())
				err := errors.NewInternal(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		} else {
			serviceError := services.DatasetService.Create(user.ID, req.Name)

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

func UploadDatasetScript(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
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

		datasetId := context.Param("id")
		datasetID, parseError := strconv.ParseUint(datasetId, 10, 64)

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

		dataset, getError := services.DatasetService.Get(uint(datasetID))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		datasetUploadDir := fileUploadDir + "datasets/" + datasetId + "/"
		if err := os.MkdirAll(datasetUploadDir, os.ModePerm); err != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "os.cmd.mkdir.dir.failed",
				TemplateData: map[string]interface{}{
					"Path":   datasetUploadDir,
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

		filePath := filepath.Join(datasetUploadDir, file.Filename)
		err := context.SaveUploadedFile(file, filePath)

		if err != nil {
			log.Printf(err.Error())
			err := errors.NewInternal(err.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
		}

		dataset.Path = filePath
		err = services.DatasetService.Update(dataset)

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

func UpdateDataset(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.DatasetReq

		if ok := util.BindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		dataset, getError := services.DatasetService.Get(req.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if dataset.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.DatasetService.Delete(req.ID)

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

func DeleteDataset(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.DatasetReq

		if ok := util.BindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		dataset, getError := services.DatasetService.Get(req.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if dataset.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.DatasetService.Delete(req.ID)

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
