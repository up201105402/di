package model

import (
	"gorm.io/gorm"
)

type Pipeline struct {
	gorm.Model
	UserID     uint `json:"owner"`
	User       User
	Name       string `json:"name"`
	Definition string `json:"definition"`
}
