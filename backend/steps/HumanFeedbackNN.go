package steps

import (
	"di/model"
	"di/util"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/guregu/null.v4"
)

type HumanFeedbackNN struct {
	ID             int
	PipelineID     uint
	RunID          uint
	IsFirstStep    bool
	Name           string
	Data_dir       null.String
	Models_dir     null.String
	Epochs_dir     null.String
	Epochs         null.Int
	Tr_fraction    null.String
	Val_fraction   null.String
	Train_desc     null.String
	Sampling       null.String
	Entropy_thresh null.String
	Nr_queries     null.Int
	IsOversampled  null.Bool
	Start_epoch    null.Int
	Dataset        null.String
}

func (step HumanFeedbackNN) GetID() int {
	return int(step.ID)
}

func (step HumanFeedbackNN) GetName() string {
	return step.Name
}

func (step *HumanFeedbackNN) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *HumanFeedbackNN) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.Data_dir = stepDescription.Data.StepConfig.Data_dir
	step.Models_dir = stepDescription.Data.StepConfig.Models_dir
	step.Epochs_dir = stepDescription.Data.StepConfig.Epochs_dir
	step.Epochs = stepDescription.Data.StepConfig.Epochs
	step.Tr_fraction = stepDescription.Data.StepConfig.Tr_fraction
	step.Val_fraction = stepDescription.Data.StepConfig.Val_fraction
	step.Train_desc = stepDescription.Data.StepConfig.Train_desc
	step.Sampling = stepDescription.Data.StepConfig.Sampling
	step.Entropy_thresh = stepDescription.Data.StepConfig.Entropy_thresh
	step.Nr_queries = stepDescription.Data.StepConfig.Nr_queries
	step.IsOversampled = stepDescription.Data.StepConfig.IsOversampled
	step.Start_epoch = stepDescription.Data.StepConfig.Start_epoch
	step.Dataset = stepDescription.Data.StepConfig.Dataset

	return nil
}

func (step *HumanFeedbackNN) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID
	return nil
}

func (step *HumanFeedbackNN) SetRunID(runID uint) error {
	step.RunID = runID
	return nil
}

func (step *HumanFeedbackNN) GetPipelineID() uint {
	return step.PipelineID
}

func (step *HumanFeedbackNN) GetRunID() uint {
	return step.RunID
}

func (step HumanFeedbackNN) Execute(logFile *os.File, I18n *i18n.Localizer) ([]model.HumanFeedbackQuery, error) {

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

	hitlDir, exists := os.LookupEnv("HITL_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "HITL_DIR",
			},
			PluralCount: 1,
		})

		runLogger.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	hitlDir = hitlDir + "/code/"

	var args []string
	args = append(args, hitlDir+"train.py")
	args, err := step.appendArgs(args, I18n, runLogger)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir
	// var stdout, stderr bytes.Buffer
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// cmdErr := cmd.Run()

	// if cmdErr != nil {
	// 	return cmdErr
	// }

	var rects [][]string
	var stoppedEpoch uint

	for epoch := step.Epochs.Int64; epoch >= 0; epoch-- {
		err := filepath.Walk(currentPipelineWorkDir+"epochs/"+string(epoch)+"/", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}

			matched, err := regexp.Match(`query_\d+_rects.csv`, []byte(info.Name()))
			if matched {
				stoppedEpoch = epoch
				csv, exists := util.ReadCsvFile(path + info.Name())
				if exists {
					for _, row := range csv {
						rects = append(rects, row)
					}
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

	var feedback []model.HumanFeedbackQuery

	for _, rect := range rects {
		feedback = append(feedback, model.HumanFeedbackQuery{
			Epoch:           uint(stoppedEpoch),
			StepID:          step.ID,
			RunID:           step.RunID,
			RectCoordinates: strings.Join(rect, ","),
		})
	}

	return feedback, nil
}

func (step HumanFeedbackNN) appendArgs(args []string, I18n *i18n.Localizer, runLogger *log.Logger) ([]string, error) {

	if step.Data_dir.Valid {
		args = append(args, "--data_dir")
		args = append(args, step.Data_dir.String)
	}
	if step.Models_dir.Valid {
		args = append(args, "--models_dir")
		args = append(args, step.Models_dir.String)
	}
	if step.Epochs.Valid {
		args = append(args, "--epochs")
		args = append(args, fmt.Sprintf("%d", step.Epochs.Int64))
	}
	if step.Tr_fraction.Valid {
		args = append(args, "--tr_fraction")
		args = append(args, step.Tr_fraction.String)
	}
	if step.Val_fraction.Valid {
		args = append(args, "--val_fraction")
		args = append(args, step.Val_fraction.String)
	}
	if step.Train_desc.Valid {
		args = append(args, "--train_desc")
		args = append(args, step.Train_desc.String)
	}
	if step.Sampling.Valid {
		args = append(args, "--sampling")
		args = append(args, step.Sampling.String)
	}
	if step.Entropy_thresh.Valid {
		args = append(args, "--entropy_thresh")
		args = append(args, step.Entropy_thresh.String)
	}
	if step.Nr_queries.Valid {
		args = append(args, "--nr_queries")
		args = append(args, fmt.Sprintf("%d", step.Nr_queries.Int64))
	}
	if step.IsOversampled.Valid {
		args = append(args, "--isOversampled")
		if step.IsOversampled.Bool {
			args = append(args, "True")
		} else {
			args = append(args, "False")
		}

	}
	if step.Start_epoch.Valid {
		args = append(args, "--start_epoch")
		args = append(args, fmt.Sprintf("%d", step.Start_epoch.Int64))
	}
	if step.Dataset.Valid {
		args = append(args, "--dataset")
		args = append(args, step.Dataset.String)
	}

	return args, nil
}
