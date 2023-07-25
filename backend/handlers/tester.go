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

func GetTesters(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		testers, getError := services.TesterService.GetByOwner(user.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		fileUploadDir := os.Getenv("FILE_UPLOAD_DIR")

		for index, tester := range testers {
			if tester.Path != "" {
				path := "/files/" + strings.Split(tester.Path, fileUploadDir)[1]
				testers[index].Path = path
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"testers": testers,
		})
	}
}

func GetTester(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		testerID := context.Param("id")

		id, parseError := strconv.ParseUint(testerID, 10, 64)

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

		tester, getError := services.TesterService.Get(uint(id))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"tester": tester,
		})
	}
}

func CreateTester(services *service.Services) gin.HandlerFunc {
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
			tester, err := services.TesterService.Get(req.ID)

			if err != nil {
				err := errors.NewNotFound(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			tester.Path = req.Path
			err = services.TesterService.Update(tester)

			if err != nil {
				log.Printf(err.Error())
				err := errors.NewInternal(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		} else {
			serviceError := services.TesterService.Create(user.ID, req.Name)

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

func UploadTesterScript(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
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

		testerId := context.Param("id")
		testerID, parseError := strconv.ParseUint(testerId, 10, 64)

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

		tester, getError := services.TesterService.Get(uint(testerID))

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		testerUploadDir := fileUploadDir + "tester/" + testerId + "/"
		if err := os.MkdirAll(testerUploadDir, os.ModePerm); err != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "os.cmd.mkdir.dir.failed",
				TemplateData: map[string]interface{}{
					"Path":   testerUploadDir,
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

		filePath := filepath.Join(testerUploadDir, file.Filename)
		err := context.SaveUploadedFile(file, filePath)

		if err != nil {
			log.Printf(err.Error())
			err := errors.NewInternal(err.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
		}

		tester.Path = filePath
		err = services.TesterService.Update(tester)

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

func UpdateTainer(services *service.Services) gin.HandlerFunc {
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

		tester, getError := services.TesterService.Get(req.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if tester.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.TesterService.Delete(req.ID)

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

func DeleteTester(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
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

		tester, getError := services.TesterService.Get(req.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewNotFound(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if tester.UserID != user.ID {
			errorMessage := fmt.Sprint("Pipeline with id %v does not belong to user %v\n", req.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewAuthorization("")
			context.JSON(err.Status(), gin.H{
				"error": errorMessage,
			})
			return
		}

		deleteError := services.TesterService.Delete(req.ID)

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
