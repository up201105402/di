package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex" json:"username"`
	Email         string `json:"email"`
	Notifications bool   `json:"notifications"`
	Avatar        string `json:"avatar"`
	Password      string `json:"password"`
}
