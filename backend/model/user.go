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
	Username string `json:"username" binding:"required,gte=5,lte=30"`
	Password string `json:"password" binding:"required,gte=5,lte=30"`
}
