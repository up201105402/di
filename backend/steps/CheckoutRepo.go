package steps

import (
	"di/model"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CheckoutRepo struct {
	ID          int
	PipelineID  uint
	RunID       uint
	IsFirstStep bool
	Name        string
	RepoURL     string `json:"repoURL"`
}

func (step CheckoutRepo) GetID() int {
	return int(step.ID)
}

func (step CheckoutRepo) GetName() string {
	return step.Name
}

func (step *CheckoutRepo) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *CheckoutRepo) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.RepoURL = stepDescription.Data.StepConfig.RepoURL.String

	return nil
}

func (step *CheckoutRepo) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *CheckoutRepo) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *CheckoutRepo) GetPipelineID() uint {
	return step.PipelineID
}

func (step *CheckoutRepo) GetRunID() uint {
	return step.RunID
}

func (step CheckoutRepo) Execute(logFile *os.File, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

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

		runLogger.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"

	if _, err := git.PlainClone(currentPipelineWorkDir, false, &git.CloneOptions{
		URL:      step.RepoURL,
		Progress: logFile,
	}); err != nil {
		// if err == git.ErrRepositoryAlreadyExists {
		// 	return err
		// }

		return nil, err
	}

	return nil, nil
}
