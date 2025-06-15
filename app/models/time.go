package models

import (
	"time"

	"github.com/nechitast/olap-backend/app/configs/clients"
)

type Dim_Time struct {
    Id_time  time.Time `json:"id" gorm:"primaryKey;type:timestamp"`
    Tahun    int       `json:"tahun" gorm:"type:integer;not null"`
    Semester int       `json:"semester" gorm:"type:integer;check:semester IN (1,2)"`
    Kuartal  string    `json:"kuartal" gorm:"type:varchar(2);check:kuartal IN ('1','2','3','4')"`
    Bulan    string    `json:"bulan" gorm:"type:varchar(20)"`
    Minggu   int       `json:"minggu" gorm:"type:integer;not null"`
    Hari     string    `json:"hari" gorm:"type:varchar(20)"`
}
func (data Dim_Time) Add() error {
	return clients.DATABASE.Model(&data).Create(&data).Error
}

func (Dim_Time) TableName() string {
	return "dim_time"
}
