package steps

import (
	"di/model"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Step interface {
	GetID() int
	GetName() string
	Execute(logFile *os.File, I18n *i18n.Localizer) ([]model.HumanFeedbackQuery, error)
	SetData(stepDescription model.NodeDescription) error
	SetPipelineID(pipelineID uint) error
	SetRunID(runID uint) error
	GetPipelineID() uint
	GetRunID() uint
	GetIsFirstStep() bool
}

type Edge interface {
	SetData(stepDescription model.NodeDescription)
	GetSourceID() int
	GetTargetID() int
}
