package model

import (
	"gorm.io/gorm"
)

type Pipeline struct {
	gorm.Model
	UserID     uint `json:"userId"`
	User       User
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

// it is used for validation and json marshalling
// it is used for validation and json marshalling
type CreatePipelineReq struct {
	User       string `json:"user"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type PipelineReq struct {
	ID         uint   `json:"id"`
	User       string `json:"user"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}
