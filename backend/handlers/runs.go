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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"pipelines": pipelines,
		})
	}
}

func FindRunsByPipeline(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		pipelineID, parseError := strconv.ParseUint(id, 10, 64)

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

		pipeline, serviceError := services.PipelineService.Get(uint(pipelineID))

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		runs, getError := services.RunService.GetByPipeline(pipeline.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"runs": runs,
		})
	}
}

func CreateRun(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		pipelineID, parseError := strconv.ParseUint(id, 10, 64)

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

		var req model.CreateRunReq

		if ok := util.BindData(context, &req); !ok {
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
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if pipeline.User.ID != user.ID {
			errorMessage := fmt.Sprint("Failed to create run for pipeline %d with user: %v\n", pipeline.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewInternal(errorMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		_, serviceError = services.RunService.Create(*pipeline)

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if req.Execute {
			serviceError := services.RunService.Execute(pipeline.ID)

			if serviceError != nil {
				log.Printf(serviceError.Error())
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

func ExecuteRun(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("runID")

		runID, parseError := strconv.ParseUint(id, 10, 64)

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

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		run, serviceError := services.RunService.Get(uint(runID))

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if run.Pipeline.User.ID != user.ID {
			errorMessage := fmt.Sprint("Failed to execute pipeline %d with user: %v\n", run.Pipeline.ID, user.Username)
			log.Printf(errorMessage)
			err := errors.NewInternal(errorMessage)
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		serviceError = services.RunService.Execute(run.ID)

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}

func FindRunResulstById(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		runID, parseError := strconv.ParseUint(id, 10, 64)

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

		run, serviceError := services.RunService.Get(uint(runID))

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		runStepStatuses, getError := services.RunService.FindRunStepStatusesByRun(run.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		runLogsDir, exists := os.LookupEnv("RUN_LOGS_DIR")

		if !exists {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "env.variable.find.failed",
				TemplateData: map[string]interface{}{
					"Name": "RUN_LOGS_DIR",
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

		logFileName, exists := os.LookupEnv("RUN_LOG_FILE_NAME")

		if !exists {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "env.variable.find.failed",
				TemplateData: map[string]interface{}{
					"Name": "RUN_LOG_FILE_NAME",
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

		runLogDir := runLogsDir + "/pipelines/" + fmt.Sprint(run.PipelineID) + "/" + fmt.Sprint(run.ID) + "/"
		logFile, err := os.ReadFile(runLogDir + logFileName)

		if err != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "os.cmd.read.file.failed",
				TemplateData: map[string]interface{}{
					"Path":   runLogDir + logFileName,
					"Reason": err.Error(),
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

		context.JSON(http.StatusOK, gin.H{
			"run":             run,
			"runStepStatuses": runStepStatuses,
			"log":             string(logFile),
		})
	}
}

func FindRunFeedbackQueriesById(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		runID, parseError := strconv.ParseUint(id, 10, 64)

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

		run, serviceError := services.RunService.Get(uint(runID))

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if run.RunStatusID != 5 {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.handler.feedback.status.error",
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

		runStepStatuses, getError := services.RunService.FindRunStepStatusesByRun(run.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		runStepStatuses = util.Filter(runStepStatuses, func(runStateStatus model.RunStepStatus) bool {
			return runStateStatus.RunStatusID == 5
		})

		if len(runStepStatuses) == 0 {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.handler.feedback.status.error",
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

		humanFeedbackQueries, err := services.RunService.FindHumanFeedbackQueriesByStepID(uint(runStepStatuses[0].StepID))

		if err != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.handler.feedback.find.fail",
				TemplateData: map[string]interface{}{
					"ID":     runStepStatuses[0].StepID,
					"Reason": err.Error(),
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

		var completeFeedbackResponse []model.HumanFeedbackQueryResponse

		for _, humanFeedbackQuery := range humanFeedbackQueries {
			feedbackRects, err := services.RunService.FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQuery.ID)

			if err != nil {
				errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.handler.feedback.find.fail",
					TemplateData: map[string]interface{}{
						"ID":     runStepStatuses[0].StepID,
						"Reason": err.Error(),
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

			pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
			currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(run.PipelineID) + "/" + fmt.Sprint(run.ID) + "/"
			epochDir := currentPipelineWorkDir + "epochs/" + fmt.Sprint(humanFeedbackQuery.Epoch) + "/"

			imagePath := epochDir + "query_" + fmt.Sprint(humanFeedbackQuery.QueryID) + "_image.png"
			_, err = os.Stat(imagePath)

			if err != nil {
				errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.handler.feedback.find.fail",
					TemplateData: map[string]interface{}{
						"QueryID": humanFeedbackQuery.QueryID,
						"Epoch":   humanFeedbackQuery.Epoch,
						"StepID":  runStepStatuses[0].StepID,
						"Reason":  err.Error(),
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

			imageURL := "/work/" + fmt.Sprint(run.PipelineID) + "/" + fmt.Sprint(run.ID) + "/epochs/" + fmt.Sprint(humanFeedbackQuery.Epoch) + "/"
			imageURL = imageURL + "query_" + fmt.Sprint(humanFeedbackQuery.QueryID) + "_image.png"

			completeFeedbackResponse = append(completeFeedbackResponse,
				model.HumanFeedbackQueryResponse{
					RunStepStatus:      runStepStatuses[0],
					HumanFeedbackQuery: humanFeedbackQuery,
					HumanFeedbackRects: feedbackRects,
					ImageURL:           imageURL,
				})
		}

		context.JSON(http.StatusOK, gin.H{
			"queries": completeFeedbackResponse,
		})
	}
}

func SubmitRunFeedback(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		id := context.Param("id")

		runID, parseError := strconv.ParseUint(id, 10, 64)

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

		run, serviceError := services.RunService.Get(uint(runID))

		if serviceError != nil {
			log.Printf(serviceError.Error())
			err := errors.NewInternal(serviceError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		if run.RunStatusID != 5 {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.handler.feedback.status.error",
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

		runStepStatuses, getError := services.RunService.FindRunStepStatusesByRun(run.ID)

		if getError != nil {
			log.Printf(getError.Error())
			err := errors.NewInternal(getError.Error())
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			return
		}

		runStepStatuses = util.Filter(runStepStatuses, func(runStateStatus model.RunStepStatus) bool {
			return runStateStatus.RunStatusID == 5
		})

		if len(runStepStatuses) == 0 {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.handler.feedback.status.error",
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

		humanFeedbackQueries, err := services.RunService.FindHumanFeedbackQueriesByStepID(uint(runStepStatuses[0].StepID))

		if err != nil {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.handler.feedback.find.fail",
				TemplateData: map[string]interface{}{
					"ID":     runStepStatuses[0].StepID,
					"Reason": err.Error(),
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

		var req model.HumanFeedbackQueryReq

		if ok := util.BindData(context, &req); !ok {
			return
		}

		for _, humanFeedbackQuery := range humanFeedbackQueries {
			feedbackRects, err := services.RunService.FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQuery.ID)

			if err != nil {
				errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.handler.feedback.find.fail",
					TemplateData: map[string]interface{}{
						"ID":     runStepStatuses[0].StepID,
						"Reason": err.Error(),
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

			for _, humanFeedbackQueryReq := range req.SingleHumanFeedbackQueryReqs {
				if humanFeedbackQueryReq.HumanFeedbackQueryID == humanFeedbackQuery.ID {
					humanFeedbackQuery.QueryStatusID = 2

					for _, feedbackRectReq := range humanFeedbackQueryReq.Rects {
						for _, feedbackRect := range feedbackRects {
							if feedbackRect.ID == feedbackRectReq.RectID {
								feedbackRect.Selected = feedbackRectReq.Selected
							}
						}
					}
				}
			}

			err = services.RunService.UpdateHumanFeedbackRects(feedbackRects)

			if err != nil {
				log.Printf(err.Error())
				err := errors.NewInternal(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}

			humanFeedbackQuery.QueryStatusID = 2

			err = services.RunService.UpdateHumanFeedbackQuery(&humanFeedbackQuery)

			if err != nil {
				log.Printf(err.Error())
				err := errors.NewInternal(err.Error())
				context.JSON(err.Status(), gin.H{
					"error": err.Message,
				})
				return
			}
		}

		context.JSON(http.StatusOK, gin.H{})
	}
}
