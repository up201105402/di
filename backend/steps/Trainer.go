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

type Trainer struct {
	ID              int
	PipelineID      uint
	RunID           uint
	IsFirstStep     bool
	Name            string
	TrainerName     string
	TrainerID       uint
	Filepath        string
	Epochs          null.Int
	IsStaggered     bool
	CustomArguments null.String
}

func (step Trainer) GetID() int {
	return int(step.ID)
}

func (step Trainer) GetName() string {
	return step.Name
}

func (step *Trainer) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *Trainer) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.Data.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.TrainerName = stepDescription.Data.StepConfig.TrainerName
	step.TrainerID = stepDescription.Data.StepConfig.TrainerID
	step.Filepath = stepDescription.Data.StepConfig.TrainerPath
	step.Epochs = stepDescription.Data.StepConfig.Epochs
	step.Filepath = stepDescription.Data.StepConfig.TrainerPath
	step.IsStaggered = stepDescription.Data.StepConfig.IsStaggered
	step.CustomArguments = stepDescription.Data.StepConfig.CustomArguments

	return nil
}

func (step *Trainer) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *Trainer) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *Trainer) GetPipelineID() uint {
	return step.PipelineID
}

func (step *Trainer) GetRunID() uint {
	return step.RunID
}

func (step *Trainer) GetIsStaggered() bool {
	return true
}

func (step Trainer) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

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

	err := step.copyOriginalTrainFile(currentPipelineWorkDir)

	if err != nil {
		return nil, err
	}

	var args []string

	args = append(args, filepath.Base(step.Filepath))

	if step.CustomArguments.Valid {
		args = append(args, step.CustomArguments.String)
	}

	var epochNumber null.Int

	if step.IsStaggered {
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
	}

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	cmdErr := cmd.Run()

	if cmdErr != nil {
		return nil, cmdErr
	}

	return step.getCreatedFeedbackQueries(epochNumber, currentPipelineWorkDir)
}

func (step Trainer) copyOriginalTrainFile(currentPipelineWorkDir string) error {
	sourceFile, err := os.Open(step.Filepath)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	defer sourceFile.Close()

	destinationFilePath := filepath.Join(currentPipelineWorkDir, filepath.Base(step.Filepath))
	trainerFileDestination, err := os.Create(destinationFilePath)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	defer trainerFileDestination.Close()

	_, err = io.Copy(trainerFileDestination, sourceFile)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	return nil
}

func (step Trainer) getCreatedFeedbackQueries(oldResumeEpoch null.Int, currentPipelineWorkDir string) ([]model.HumanFeedbackQueryPayload, error) {
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
