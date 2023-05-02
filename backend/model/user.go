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
	Password      string `json:"-"`
}

// it is used for validation and json marshalling
type UserReq struct {
	Username string `json:"username" binding:"required,gte=5,lte=30"`
	Password string `json:"password" binding:"required,gte=5,lte=30"`
}
