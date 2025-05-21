package handlers

import (
	"net/http"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nechitast/olap-backend/app/models"
	"github.com/nechitast/olap-backend/app/configs/clients"
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

//fungsi get hotspot
func GetHotspot(ctx *fiber.Ctx) error {
	var hotspot []models.Fact_Hotspot

	fmt.Println("Querying database for hotspots...")

	if err := clients.DATABASE.
		Preload("Dim_Location").
		Preload("Dim_Time").
		Preload("Dim_Confidence").
		Preload("Dim_Satelite").
		Find(&hotspot).Error; err != nil {
		fmt.Println("Database error:", err)
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	fmt.Println("Hotspots retrieved:", len(hotspot))
	if len(hotspot) == 0 {
		fmt.Println("No data found in fact_hotspot!")
	}

	// Format ke GeoJSON
	geoJSON := fiber.Map{
		"type":     "FeatureCollection",
		"features": []fiber.Map{},
	}

	for _, h := range hotspot {
		longitude, latitude, err := h.Dim_Location.ExtractLatLng()
		if err != nil {
			fmt.Printf("Error extracting lat/lng: %v\n", err)
			continue
		}

		fmt.Printf("Hotspot ID: %d, Location: (%f, %f)\n", h.ID_Location, longitude, latitude)

		geoJSON["features"] = append(geoJSON["features"].([]fiber.Map), fiber.Map{
			"type": "Feature",
			"properties": fiber.Map{
				"confidence":    h.Dim_Confidence.Confidence_Level,
				"satellite":     h.Dim_Satelite.Satelite_Name,
				"time":          h.Dim_Time.Id_time,
				"hotspot_count": h.Hotspot_Count,
				"location": fiber.Map {
					"pulau": h.Dim_Location.Pulau,
					"provinsi": h.Dim_Location.Provinsi,
					"kab_kota": h.Dim_Location.Kab_kota,
					"kecamatan": h.Dim_Location.Kecamatan,
					"desa": h.Dim_Location.Desa,
					"geom_desa": h.Dim_Location.GeomDesa,
				},
			},
			"geometry": fiber.Map{
				"type":        "Point",
				"coordinates": []float64{latitude, longitude},
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
	var locations []models.Dim_Location

	// Ambil data dari database
	if err := clients.DATABASE.Find(&locations).Error; err != nil {
		return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
	}

	// Format data untuk frontend
	formattedLocations := make([]map[string]interface{}, 0)
	for _, loc := range locations {
		lat, lng, err := loc.ExtractLatLng() 
		if err != nil {
			fmt.Printf("Error extracting lat/lng for location %d: %v\n", loc.Id_location, err)
			continue
		}

		formattedLocations = append(formattedLocations, map[string]interface{}{
			"pulau":       loc.Pulau,
			"provinsi":    loc.Provinsi,
			"kab_kota":	   loc.Kab_kota,
			"kecamatan":   loc.Kecamatan,
			"desa":		   loc.Desa,
			"lat":         lat,
			"lng":         lng,

		})
	}

	return ResponseJson(ctx, http.StatusOK, formattedLocations)
}
