package model

import (
	"time"

	"gorm.io/gorm"
)

type RunStatus struct {
	gorm.Model
	Name        string
	IsFinal     bool
	Description string
}

type Run struct {
	gorm.Model
	PipelineID          uint `json:"pipelineId"`
	Pipeline            Pipeline
	RunStatusID         uint
	RunStatus           RunStatus
	ErrorMessage        string
	Definition          string
	StepWaitingFeedback int
	LastRun             time.Time
}

type RunStepStatus struct {
	gorm.Model
	Name         string
	StepID       int
	RunID        uint
	Run          Run
	RunStatusID  uint
	RunStatus    RunStatus
	ErrorMessage string
	LastRun      time.Time
}

type HumanFeedbackQueryPayload struct {
	Epoch   uint
	StepID  int
	QueryID uint
	RunID   uint
	Rects   [][]uint
}

type HumanFeedbackQuery struct {
	gorm.Model
	Epoch         uint `gorm:"index:idx_member"`
	StepID        int  `gorm:"index:idx_member"`
	QueryID       uint `gorm:"index:idx_member"`
	RunID         uint `gorm:"index:idx_member"`
	Run           Run
	QueryStatusID uint
	QueryStatus   QueryStatus
}

type QueryStatus struct {
	gorm.Model
	Name    string
	IsFinal bool
}

type HumanFeedbackRect struct {
	gorm.Model
	X1                   uint
	Y1                   uint
	X2                   uint
	Y2                   uint
	HumanFeedbackQueryID uint `gorm:"index"`
	HumanFeedbackQuery   HumanFeedbackQuery
	Selected             bool
}

type HumanFeedbackQueryResponse struct {
	RunStepStatus      RunStepStatus
	HumanFeedbackQuery HumanFeedbackQuery
	HumanFeedbackRects []HumanFeedbackRect
	ImageURL           string
}

type CreateRunReq struct {
	Execute bool `json:"execute"`
}

type ExecuteRunReq struct {
	ID uint `json:"id"`
}
