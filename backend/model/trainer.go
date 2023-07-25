package model

import "gorm.io/gorm"

type Trainer struct {
	gorm.Model
	UserID uint `json:"userId"`
	User   User
	Name   string `json:"name" gorm:"uniqueIndex"`
	Path   string `json:"path"`
}

type TrainerReq struct {
	ID   uint   `json:"id"`
	User string `json:"user"`
	Name string `json:"name"`
	Path string `json:"path"`
}
