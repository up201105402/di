package model

import (
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
	PipelineID  uint `json:"pipelineId"`
	Pipeline    Pipeline
	Status      RunStatus
	Description string
}

type CreateRunReq struct {
	ID         uint `json:"id"`
	PipelineID uint `json:"pipelineId"`
	Execute    bool `json:"execute"`
}

type ExecuteRunReq struct {
	ID uint `json:"id"`
}
