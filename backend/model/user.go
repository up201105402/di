package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Language string `json:"language"`
}

type UserReq struct {
	Username    string `json:"username" binding:"gte=5,lte=30"`
	Password    string `json:"password" binding:"gte=5,lte=30"`
	NewPassword string `json:"newPassword"`
}
