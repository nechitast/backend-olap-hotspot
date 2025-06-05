package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nechitast/olap-backend/app/configs/clients"
	"github.com/nechitast/olap-backend/app/models"
)

func AddHotspot(ctx *fiber.Ctx) error {
	var data models.Fact_Hotspot

	if err := ctx.BodyParser(&data); err != nil {
		return ResponseJson(ctx, http.StatusBadRequest, err.Error())
	}

	if err := data.Add(); err != nil {
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	return ResponseJson(ctx, http.StatusOK, data)
}

func GetHotspot(ctx *fiber.Ctx) error {
	var hotspotsData []struct {
		Hotspot_Count int     `gorm:"column:hotspot_count"`
		Hotspot_Time  string  `gorm:"column:hotspot_time"`
		Desa          string  `gorm:"column:desa"`
		Kecamatan     string  `gorm:"column:kecamatan"`
		Kab_kota      string  `gorm:"column:kab_kota"`
		Provinsi      string  `gorm:"column:provinsi"`
		Pulau         string  `gorm:"column:pulau"`
		Longitude     float64 `gorm:"column:longitude"`
		Latitude      float64 `gorm:"column:latitude"`
		Confidence    string  `gorm:"column:confidence_level"`
		Satelite      string  `gorm:"column:satelite_name"`
		Time          string  `gorm:"column:id_time"`
		Minggu        string  `gorm:"column:minggu"`
	}

	fmt.Println("Querying database for hotspots...")

	err := clients.DATABASE.Raw(`
		SELECT
			fh.hotspot_count,
			fh.hotspot_time::text,
			dl.desa,
			dl.kecamatan,
			dl.kab_kota,
			dl.provinsi,
			dl.pulau,
			ST_X(dl.geom_desa) AS longitude,
			ST_Y(dl.geom_desa) AS latitude,
			dc.confidence_level,
			ds.satelite_name,
			dt.id_time::text,
			dt.minggu::text
		FROM fact_hotspot fh
		JOIN dim_location dl ON fh.id_location = dl.id_location
		JOIN dim_confidence dc ON fh.id_confidence = dc.id_confidence
		JOIN dim_satelite ds ON fh.id_satelite = ds.id_satelite
		JOIN dim_time dt ON fh.id_time = dt.id_time
	`).Scan(&hotspotsData).Error

	if err != nil {
		fmt.Println("Database error in GetHotspot:", err)
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	fmt.Println("Hotspots retrieved:", len(hotspotsData))
	if len(hotspotsData) == 0 {
		fmt.Println("No data found in fact_hotspot for this query!")
	}

	// Format ke GeoJSON
	geoJSON := fiber.Map{
		"type":     "FeatureCollection",
		"features": []fiber.Map{},
	}

	for _, h := range hotspotsData {
		geoJSON["features"] = append(geoJSON["features"].([]fiber.Map), fiber.Map{
			"type": "Feature",
			"properties": fiber.Map{
				"confidence":    h.Confidence,
				"satellite":     h.Satelite,
				"time":          h.Time,
				"minggu":		 h.Minggu,
				"hotspot_count": h.Hotspot_Count,
				"hotspot_time":  h.Hotspot_Time,
				"location": fiber.Map{
					"pulau":     h.Pulau,
					"provinsi":  h.Provinsi,
					"kab_kota":  h.Kab_kota,
					"kecamatan": h.Kecamatan,
					"desa":      h.Desa,
				},
			},
			"geometry": fiber.Map{
				"type":        "Point",
				"coordinates": []float64{h.Longitude, h.Latitude},
			},
		})
	}

	return ctx.JSON(geoJSON)
}

func AddLocation(ctx *fiber.Ctx) error {
	var data models.Dim_Location

	if err := ctx.BodyParser(&data); err != nil {
		return ResponseJson(ctx, http.StatusBadRequest, err.Error())
	}

	if err := data.Add(); err != nil {
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	return ResponseJson(ctx, http.StatusOK, data)
}

func AddTime(ctx *fiber.Ctx) error {
	var data models.Dim_Time

	if err := ctx.BodyParser(&data); err != nil {
		return ResponseJson(ctx, http.StatusBadRequest, err.Error())
	}

	if err := data.Add(); err != nil {
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	return ResponseJson(ctx, http.StatusOK, data)
}

func AddConfidence(ctx *fiber.Ctx) error {
	var data models.Dim_Confidence

	if err := ctx.BodyParser(&data); err != nil {
		return ResponseJson(ctx, http.StatusBadRequest, err.Error())
	}

	if err := data.Add(); err != nil {
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	return ResponseJson(ctx, http.StatusOK, data)
}

func AddSatelite(ctx *fiber.Ctx) error {
	var data models.Dim_Satelite

	if err := ctx.BodyParser(&data); err != nil {
		return ResponseJson(ctx, http.StatusBadRequest, err.Error())
	}

	if err := data.Add(); err != nil {
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	return ResponseJson(ctx, http.StatusOK, data)
}

func QueryLocation(ctx *fiber.Ctx) error {
	var locations []struct {
		Id_location int    `gorm:"column:id_location"`
		Pulau       string `gorm:"column:pulau"`
		Provinsi    string `gorm:"column:provinsi"`
		Kab_kota    string `gorm:"column:kab_kota"`
		Kecamatan   string `gorm:"column:kecamatan"`
		Desa        string `gorm:"column:desa"`
		Longitude   float64 `gorm:"column:longitude"` 
		Latitude    float64 `gorm:"column:latitude"` 
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
		fmt.Printf("Database error in QueryLocation: %v\n", err)
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	fmt.Printf("Successfully fetched %d locations from DB in QueryLocation.\n", len(locations))

	formattedLocations := make([]map[string]interface{}, 0)
	for _, loc := range locations {
		formattedLocations = append(formattedLocations, map[string]interface{}{
			"id": loc.Id_location,
			"pulau":     loc.Pulau,
			"provinsi":  loc.Provinsi,
			"kab_kota":  loc.Kab_kota,
			"kecamatan": loc.Kecamatan,
			"desa":      loc.Desa,
			"lat":       loc.Latitude,
			"lng":       loc.Longitude,
		})
	}

	return ResponseJson(ctx, http.StatusOK, formattedLocations)
}