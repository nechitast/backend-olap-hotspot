package models

import (
	"fmt"
	"errors"
	"github.com/nechitast/olap-backend/app/configs/clients"
)

type Dim_Location struct {
	Id_location int    `json:"id" gorm:"primaryKey"`
	Pulau string `json:"pulau" form:"pulau"`
	Provinsi string `json:"provinsi" form:"provinsi"`
	Kab_kota string `json:"kab_kota" form:"kab_kota"`
	Kecamatan string `json:"kecamatan" form:"kecamatan"`
	Desa string `json:"desa" form:"desa"`
	GeomDesa []byte `json:"geom_desa" gorm:"column:geom_desa"` 
}

func (d *Dim_Location) ExtractLatLng() (float64, float64, error) {
    return 0, 0, errors.New("ExtractLatLng is causing WKB unmarshal errors and is not used by GetAllLocations anymore.")
}

func (data Dim_Location) Add() error {
	return clients.DATABASE.Model(&data).Create(&data).Error
}

func GetAllLocations() ([]map[string]interface{}, error) {
	var locations []struct {
		Id_location int     `gorm:"column:id_location"`
		Pulau string  `gorm:"column:pulau"`
		Provinsi string  `gorm:"column:provinsi"`
		Kab_kota string  `gorm:"column:kab_kota"`
		Kecamatan string  `gorm:"column:kecamatan"`
		Desa string  `gorm:"column:desa"`
		Longitude float64 `gorm:"column:longitude"` 
		Latitude float64 `gorm:"column:latitude"` 
	}

	err := clients.DATABASE.Raw(`
		SELECT 
			id_location, 
			pulau, 
			provinsi, 
			kab_kota, 
			kecamatan, 
			desa, 
			ST_X(geom_desa) as longitude,
			ST_Y(geom_desa) as latitude
		FROM dim_location
	`).Scan(&locations).Error 

	if err != nil {
		fmt.Printf("Database error fetching locations from DB: %v\n", err) 
		return nil, err
	}

	fmt.Printf("Successfully fetched %d locations from DB.\n", len(locations)) 

	var result []map[string]interface{}
	for _, loc := range locations {
		result = append(result, map[string]interface{}{
			"id": loc.Id_location,
			"pulau": loc.Pulau,
			"provinsi": loc.Provinsi,
			"kab_kota": loc.Kab_kota,
			"kecamatan": loc.Kecamatan,
			"desa": loc.Desa,
			"longitude": loc.Longitude,
			"latitude": loc.Latitude,
		})
	}

	return result, nil
}

func (Dim_Location) TableName() string {
	return "dim_location"
}