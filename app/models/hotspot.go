package models

import (
	"time"
	"github.com/nechitast/olap-backend/app/configs/clients"
)

type Fact_Hotspot struct {
	ID_Location   int       `json:"-" form:"location"`
	ID_Confidence int       `json:"-" form:"confidence"`
	ID_Time       time.Time `json:"-" form:"time"`
	ID_Satelite   int       `json:"-" form:"satelite"`
	Hotspot_Count int       `json:"hotspot_count" form:"total"`

	Dim_Location   Dim_Location   `json:"dim_location" gorm:"foreignKey:ID_Location"`
	Dim_Time       Dim_Time       `json:"dim_time" gorm:"foreignKey:ID_Time"`
	Dim_Confidence Dim_Confidence `json:"dim_confidence" gorm:"foreignKey:ID_Confidence"`
	Dim_Satelite   Dim_Satelite   `json:"dim_satelite" gorm:"foreignKey:ID_Satelite"`
}

func (data Fact_Hotspot) Add() error {
	return clients.DATABASE.Model(&data).Omit("Location", "Time", "Confidence", "Satelite").Create(&data).Error
}

func (Fact_Hotspot) TableName() string {
    return "fact_hotspot"
}
