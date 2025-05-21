package models

import (
	"github.com/nechitast/olap-backend/app/configs/clients"
	"fmt"
	"errors"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

type Dim_Location struct {
	Id_location int    `json:"id" gorm:"primaryKey"`
	Pulau       string `json:"pulau" form:"pulau"`
	Provinsi    string `json:"provinsi" form:"provinsi"`
	Kab_kota    string `json:"kab_kota" form:"kab_kota"`
	Kecamatan   string `json:"kecamatan" form:"kecamatan"`
	Desa        string `json:"desa" form:"desa"`
	GeomDesa    []byte `json:"geom_desa" gorm:"column:geom_desa"` 
}

func (d *Dim_Location) ExtractLatLng() (float64, float64, error) {
	if len(d.GeomDesa) == 0 {
		return 0, 0, errors.New("GeomDesa is empty")
	}

	geom, err := wkb.Unmarshal(d.GeomDesa)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal WKB: %v", err)
	}

	point, ok := geom.(orb.Point)
	if !ok {
		return 0, 0, errors.New("Geom_Desa is not a POINT type")
	}

	return point[1], point[0], nil
}

func (data Dim_Location) Add() error {
	return clients.DATABASE.Model(&data).Create(&data).Error
}

func GetAllLocations() ([]map[string]interface{}, error) {
	var locations []Dim_Location
	err := clients.DATABASE.Select("id_location, pulau, provinsi, kab_kota, kecamatan, desa, ST_AsBinary(geom_desa) as geom_desa").Find(&locations).Error
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, loc := range locations {
		lng, lat, err := loc.ExtractLatLng()
		if err != nil {
			fmt.Printf("Error extracting lat/lng: %v\n", err)
			continue
		}

		result = append(result, map[string]interface{}{
			"id":        loc.Id_location,
			"pulau":     loc.Pulau,
			"provinsi":  loc.Provinsi,
			"kab_kota":  loc.Kab_kota,
			"kecamatan": loc.Kecamatan,
			"desa":      loc.Desa,
			"longitude": lng,
			"latitude":  lat,
		})
	}

	return result, nil
}

func (Dim_Location) TableName() string {
	return "dim_location"
}
