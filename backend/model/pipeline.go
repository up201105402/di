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

type StepDescription struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Data  struct {
		NameAndType struct {
			NodeName string `json:"nodeName"`
			NodeType string `json:"nodeType"`
		} `json:"nameAndType"`
		StepConfig struct {
			RepoURL        string `json:"repoURL"`
			Directory      string `json:"trainingSetDirectory"`
			Fraction       string `json:"fraction"`
			ModelDirectory string `json:"modelDirectory"`
			Epochs         string `json:"epochs"`
		} `json:"stepConfig"`
		IsFirstStep bool `json:"isFirstStep"`
	} `json:"data"`
}

type Step interface {
	Execute() error
}
