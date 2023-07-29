package model

import "gorm.io/gorm"

type Dataset struct {
	gorm.Model
	UserID uint `json:"userId"`
	User   User
	Name   string `json:"name" gorm:"uniqueIndex"`
	Path   string `json:"path"`
}

type DatasetReq struct {
	ID   uint   `json:"id"`
	User string `json:"user"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type DatasetScript struct {
	ID        uint `json:"id"`
	DatasetID uint `json:"datasetId"`
	Dataset   Dataset
	Name      string `json:"name"`
	Path      string `json:"path"`
}
