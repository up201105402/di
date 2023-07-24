package steps

import (
	"di/model"
	"di/util"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/guregu/null.v4"
)

type CustomHITL struct {
	ID                  int
	PipelineID          uint
	RunID               uint
	IsFirstStep         bool
	Name                string
	CustomTrainFilename string
	Filename            string
	CustomArguments     null.String
	Epochs_dir          null.String
	Epochs              null.Int
	Start_epoch         null.Int
}

func (step CustomHITL) GetID() int {
	return int(step.ID)
}

func (step CustomHITL) GetName() string {
	return step.Name
}

func (step *CustomHITL) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *CustomHITL) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.Filename = stepDescription.Data.StepConfig.Filename.String
	step.CustomArguments = stepDescription.Data.StepConfig.CustomArguments
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.Epochs_dir = stepDescription.Data.StepConfig.Epochs_dir
	step.Epochs = stepDescription.Data.StepConfig.Epochs
	step.Start_epoch = stepDescription.Data.StepConfig.Start_epoch
	step.CustomTrainFilename = "custom_train.py"

	return nil
}

func (step *CustomHITL) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID
	return nil
}

func (step *CustomHITL) SetRunID(runID uint) error {
	step.RunID = runID
	return nil
}

func (step *CustomHITL) GetPipelineID() uint {
	return step.PipelineID
}

func (step *CustomHITL) GetRunID() uint {
	return step.RunID
}

func (step *CustomHITL) GetIsStaggered() bool {
	return true
}

func (step CustomHITL) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})
		return nil, errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"
	customHITLDir := currentPipelineWorkDir + "custom_hitl/"

	trainPY := filepath.Join(customHITLDir, step.CustomTrainFilename)
	var args []string
	args = append(args, trainPY)
	args, err := step.appendArgs(args, currentPipelineWorkDir, I18n, runLogger)

	if err != nil {
		return nil, err
	}

	var epochNumber null.Int

	for _, rects := range feedbackRects {
		if len(rects) > 0 {
			epochNumber = null.NewInt(int64(rects[0].HumanFeedbackQuery.Epoch), true)
			queryNumber := rects[0].HumanFeedbackQuery.QueryID
			epochDir := currentPipelineWorkDir + "epochs/" + fmt.Sprint(epochNumber.Int64) + "/"
			fileName := epochDir + "query_" + fmt.Sprint(queryNumber) + "_rects_selected.csv"
			rectsSelectedFile, err := os.Create(fileName)

			if err != nil {
				errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "os.cmd.create.file.failed",
					TemplateData: map[string]interface{}{
						"Path":   fileName,
						"Reason": err.Error(),
					},
					PluralCount: 1,
				})

				runLogger.Println(errMessage)
				return nil, errors.New(errMessage)
			}

			csvWriter := csv.NewWriter(rectsSelectedFile)

			rects = util.Filter(rects, func(rect model.HumanFeedbackRect) bool {
				return rect.Selected == true
			})

			var rectLines [][]string

			for _, rect := range rects {
				rectLines = append(rectLines, []string{fmt.Sprint(rect.X1), fmt.Sprint(rect.Y1), fmt.Sprint(rect.X2), fmt.Sprint(rect.Y2)})
			}

			err = csvWriter.WriteAll(rectLines)

			if err != nil {
				errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "os.cmd.create.file.failed",
					TemplateData: map[string]interface{}{
						"Path":   fileName,
						"Reason": err.Error(),
					},
					PluralCount: 1,
				})

				runLogger.Println(errMessage)
				return nil, errors.New(errMessage)
			}

			err = rectsSelectedFile.Close()

			if err != nil {
				errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "os.cmd.create.file.failed",
					TemplateData: map[string]interface{}{
						"Path":   fileName,
						"Reason": err.Error(),
					},
					PluralCount: 1,
				})

				runLogger.Println(errMessage)
				return nil, errors.New(errMessage)
			}
		}
	}

	if epochNumber.Valid {
		args = append(args, "--resume_epoch")
		args = append(args, fmt.Sprintf("%d", epochNumber.Int64))
	}

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir
	// var stdout, stderr bytes.Buffer
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	cmdErr := cmd.Run()

	if cmdErr != nil {
		return nil, cmdErr
	}

	return step.getCreatedFeedbackQueries(epochNumber, currentPipelineWorkDir)
}

func (step CustomHITL) appendArgs(args []string, currentPipelineWorkDir string, I18n *i18n.Localizer, runLogger *log.Logger) ([]string, error) {

	if step.Epochs.Valid {
		args = append(args, "--epochs")
		args = append(args, fmt.Sprintf("%d", step.Epochs.Int64))
	}
	if step.Start_epoch.Valid {
		args = append(args, "--start_epoch")
		args = append(args, fmt.Sprintf("%d", step.Start_epoch.Int64))
	}

	args = append(args, "--epochs_dir")
	args = append(args, currentPipelineWorkDir+"epochs/")

	return args, nil
}

func (step CustomHITL) getCreatedFeedbackQueries(oldResumeEpoch null.Int, currentPipelineWorkDir string) ([]model.HumanFeedbackQueryPayload, error) {
	var feedback []model.HumanFeedbackQueryPayload
	var stoppedEpoch uint

	for epoch := step.Epochs.Int64; epoch >= 0; epoch-- {
		if oldResumeEpoch.Valid && oldResumeEpoch.Int64 == epoch {
			break
		}

		err := filepath.Walk(currentPipelineWorkDir+"epochs/"+fmt.Sprintf("%d", epoch)+"/", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			matched, err := regexp.Match(`query_\d+_rects.csv`, []byte(info.Name()))
			if matched {
				stoppedEpoch = uint(epoch)

				rgx := regexp.MustCompile(`query_(\d+?)_rects\.csv`)
				extractedGroups := rgx.FindStringSubmatch(info.Name())
				queryID, err := strconv.ParseUint(extractedGroups[1], 10, 64)
				if err != nil {
					return err
				}

				csv, exists := util.ReadCsvFile(path)
				if exists {
					var rects [][]uint

					for _, row := range csv {
						var rect []uint

						for _, cell := range row {
							coor, err := strconv.ParseUint(cell, 10, 64)

							if err != nil {
								return err
							}

							rect = append(rect, uint(coor))
						}

						if len(rect) == 4 {
							rects = append(rects, rect)
						}
					}

					feedback = append(feedback, model.HumanFeedbackQueryPayload{
						Epoch:   uint(stoppedEpoch),
						StepID:  step.ID,
						RunID:   step.RunID,
						QueryID: uint(queryID),
						Rects:   rects,
					})
				}
			}

			return nil
		})

		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}

	return feedback, nil
}

func (step CustomHITL) createTrainFile(logFile *os.File, I18n *i18n.Localizer) error {

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"
	customHITLDir := currentPipelineWorkDir + "custom_hitl/"

	if err := os.MkdirAll(customHITLDir, os.ModePerm); err != nil {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   customHITLDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

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
		return errors.New(errMessage)
	}

	sourceFile, err := os.Open(fileUploadDir + "pipelines/" + fmt.Sprint(step.PipelineID) + "/" + step.Filename)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	defer sourceFile.Close()

	customTrainFile, err := os.Create(customHITLDir + step.CustomTrainFilename)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	defer customTrainFile.Close()

	_, err = io.Copy(customTrainFile, sourceFile)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	return nil
}
