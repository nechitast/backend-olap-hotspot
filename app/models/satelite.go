package models

import (
	"github.com/nechitast/olap-backend/app/configs/clients"
)

type Dim_Satelite struct {
	Id_satelite   int    `json:"id" gorm:"primaryKey"`
	Satelite_Name string `json:"name" form:"name"`
}

func (data Dim_Satelite) Add() error {
	return clients.DATABASE.Model(&data).Create(&data).Error
}

func (Dim_Satelite) TableName() string {
	return "dim_satelite"
}
