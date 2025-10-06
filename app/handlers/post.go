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
    // Get query parameters for filtering
    confidence := ctx.Query("confidence")
    satelite := ctx.Query("satelite")
    pulau := ctx.Query("pulau")
    provinsi := ctx.Query("provinsi")
    kota := ctx.Query("kota")
    kecamatan := ctx.Query("kecamatan")
    desa := ctx.Query("desa")
    tahun := ctx.Query("tahun")
    semester := ctx.Query("semester")
    kuartal := ctx.Query("kuartal")
    bulan := ctx.Query("bulan")
    minggu := ctx.Query("minggu")
    selectedDate := ctx.Query("selectedDate")

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
        Tahun         int     `gorm:"column:tahun"`
        Semester      string  `gorm:"column:semester"`
        Kuartal       string  `gorm:"column:kuartal"`
        Bulan         string  `gorm:"column:bulan"`
    }

    fmt.Printf("Querying hotspots with selectedDate: %s\n", selectedDate)

    limit := ctx.QueryInt("limit", 1000)
    offset := ctx.QueryInt("offset", 0)
    
    // Build dynamic query with WHERE conditions
    query := `
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
            dt.minggu::text,
            dt.tahun,
            dt.semester::text,
            dt.kuartal::text,
            dt.bulan::text
        FROM fact_hotspot fh
        JOIN dim_location dl ON fh.id_location = dl.id_location
        JOIN dim_confidence dc ON fh.id_confidence = dc.id_confidence
        JOIN dim_satelite ds ON fh.id_satelite = ds.id_satelite
        JOIN dim_time dt ON fh.id_time = dt.id_time
        WHERE 1=1`

    var args []interface{}
    argIndex := 1

    if selectedDate != "" {
        query += fmt.Sprintf(" AND DATE(fh.hotspot_time) = $%d", argIndex)
        args = append(args, selectedDate)
        argIndex++
        fmt.Printf("Applied date filter: %s\n", selectedDate)
    }

    if confidence != "" {
        query += fmt.Sprintf(" AND dc.confidence_level = $%d", argIndex)
        args = append(args, confidence)
        argIndex++
    }
    
    if satelite != "" {
        query += fmt.Sprintf(" AND ds.satelite_name = $%d", argIndex)
        args = append(args, satelite)
        argIndex++
    }

    if pulau != "" {
        query += fmt.Sprintf(" AND dl.pulau = $%d", argIndex)
        args = append(args, pulau)
        argIndex++
    }

    if provinsi != "" {
        query += fmt.Sprintf(" AND dl.provinsi = $%d", argIndex)
        args = append(args, provinsi)
        argIndex++
    }

    if kota != "" {
        query += fmt.Sprintf(" AND dl.kab_kota = $%d", argIndex)
        args = append(args, kota)
        argIndex++
    }

    if kecamatan != "" {
        query += fmt.Sprintf(" AND dl.kecamatan = $%d", argIndex)
        args = append(args, kecamatan)
        argIndex++
    }

    if desa != "" {
        query += fmt.Sprintf(" AND dl.desa = $%d", argIndex)
        args = append(args, desa)
        argIndex++
    }

    // Time filters
    if tahun != "" {
        query += fmt.Sprintf(" AND dt.tahun = $%d", argIndex)
        args = append(args, tahun)
        argIndex++
    }

    if semester != "" {
        query += fmt.Sprintf(" AND dt.semester = $%d", argIndex)
        args = append(args, semester)
        argIndex++
    }

    if kuartal != "" {
        query += fmt.Sprintf(" AND dt.kuartal = $%d", argIndex)
        args = append(args, kuartal)
        argIndex++
    }

    if bulan != "" {
        query += fmt.Sprintf(" AND dt.bulan = $%d", argIndex)
        args = append(args, bulan)
        argIndex++
    }

    if minggu != "" {
        query += fmt.Sprintf(" AND dt.minggu = $%d", argIndex)
        args = append(args, minggu)
        argIndex++
    }

    fmt.Printf("Final query: %s\n", query)
    fmt.Printf("Query args: %v\n", args)

    query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
    err := clients.DATABASE.Raw(query, args...).Scan(&hotspotsData).Error

    if err != nil {
        fmt.Printf("Database error in GetHotspot: %v\n", err)
        return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
    }

    fmt.Printf("Hotspots retrieved with filters: %d\n", len(hotspotsData))

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
                "minggu":        h.Minggu,
                "tahun":         h.Tahun,
                "semester":      h.Semester,
                "kuartal":       h.Kuartal,
                "bulan":         h.Bulan,
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

func QueryLocationAggregated(ctx *fiber.Ctx) error {
    // Get query parameters for filtering
    selectedDate := ctx.Query("selectedDate")
    confidence := ctx.Query("confidence")
    satelite := ctx.Query("satelite")
    pulau := ctx.Query("pulau")
    provinsi := ctx.Query("provinsi")
    kota := ctx.Query("kota")

    fmt.Printf("QueryLocationAggregated with selectedDate: %s\n", selectedDate)

    var locationData []struct {
        Location string `gorm:"column:location"`
        Total    int    `gorm:"column:total"`
    }

    // Build dynamic query for location aggregation with filters
    query := `
        SELECT 
            dl.pulau as location,
            COUNT(*) as total
        FROM fact_hotspot fh
        JOIN dim_location dl ON fh.id_location = dl.id_location
        JOIN dim_confidence dc ON fh.id_confidence = dc.id_confidence
        JOIN dim_satelite ds ON fh.id_satelite = ds.id_satelite
        JOIN dim_time dt ON fh.id_time = dt.id_time
        WHERE 1=1`

    var args []interface{}
    argIndex := 1

    // Apply date filter
    if selectedDate != "" {
        query += fmt.Sprintf(" AND DATE(fh.hotspot_time) = $%d", argIndex)
        args = append(args, selectedDate)
        argIndex++
    }

    // Apply other filters
    if confidence != "" {
        query += fmt.Sprintf(" AND dc.confidence_level = $%d", argIndex)
        args = append(args, confidence)
        argIndex++
    }
    
    if satelite != "" {
        query += fmt.Sprintf(" AND ds.satelite_name = $%d", argIndex)
        args = append(args, satelite)
        argIndex++
    }

    if pulau != "" {
        query += fmt.Sprintf(" AND dl.pulau = $%d", argIndex)
        args = append(args, pulau)
        argIndex++
    }

    if provinsi != "" {
        query += fmt.Sprintf(" AND dl.provinsi = $%d", argIndex)
        args = append(args, provinsi)
        argIndex++
    }

    if kota != "" {
        query += fmt.Sprintf(" AND dl.kab_kota = $%d", argIndex)
        args = append(args, kota)
        argIndex++
    }

    query += " GROUP BY dl.pulau ORDER BY total DESC"

    fmt.Printf("Location aggregation query: %s\n", query)
    fmt.Printf("Query args: %v\n", args)

    err := clients.DATABASE.Raw(query, args...).Scan(&locationData).Error
    if err != nil {
        fmt.Printf("Database error in QueryLocationAggregated: %v\n", err)
        return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
    }

    // Format as [[location, total]]
    result := make([][]interface{}, 0)
    for _, item := range locationData {
        result = append(result, []interface{}{item.Location, item.Total})
    }

    fmt.Printf("Location aggregation result: %v\n", result)
    return ctx.JSON(result)
}