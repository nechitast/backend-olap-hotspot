package models

import (
	"github.com/nechitast/olap-backend/app/configs/clients"
)

type Dim_Confidence struct {
	Id_confidence    int    `json:"id" gorm:"primaryKey"`
	Confidence_Level string `json:"level" form:"level"`
}

func (data Dim_Confidence) Add() error {
	return clients.DATABASE.Model(&data).Create(&data).Error
}

func (Dim_Confidence) TableName() string {
	return "dim_confidence"
}
