package handlers

import (
	"net/http"
	"github.com/nechitast/olap-backend/app/configs/clients"
	"github.com/gofiber/fiber/v2"
	"github.com/nechitast/olap-backend/app/helpers"
	"github.com/nechitast/olap-backend/app/models/payload"
	"github.com/xdbsoft/olap"
	"fmt"
)

func GetHeader(ctx *fiber.Ctx) error {
	var cube olap.Cube
	var err error
	var other payload.Other
	limit := ctx.QueryInt("limit", 100)
	offset := ctx.QueryInt("offset", 10)

	cube, err = helpers.SQLtoCube(other, limit, offset)
	if err != nil {
		return ResponseJson(ctx, http.StatusBadGateway, err.Error())
	}

	return ResponseJson(ctx, http.StatusOK, cube.Headers())
}

func Query(ctx *fiber.Ctx) error {
	var cube olap.Cube
	var err error

	dimension := ctx.Params("dimension")
	if dimension == "time" {
		return handleTimeQuery(ctx)
	}

	selectedDate := ctx.Query("selectedDate")
    if selectedDate != "" && dimension == "location" {
        fmt.Printf("Detected selectedDate filter: %s for dimension: %s\n", selectedDate, dimension)
    	return handleLocationQueryWithDate(ctx, selectedDate)
    }

	point := ctx.Query("point")
	dim := ctx.Query("dimension")
	limit := ctx.QueryInt("limit", 100)
	offset := ctx.QueryInt("offset", 10)

	location, time, other := helpers.GetPayload(ctx)

	if dimension == "location" || dimension == "time" {
		if location.Pulau != "" || time.Tahun != "" {
			if location.Pulau != "" && time.Tahun == "" {
				cube, err = helpers.CubeLocation(location, other)
				if err != nil {
					return ResponseJson(ctx, http.StatusBadGateway, err.Error())
				}
			}
			if location.Pulau == "" && time.Tahun != "" {
				cube, err = helpers.CubeTime(time, other)
				if err != nil {
					return ResponseJson(ctx, http.StatusBadGateway, err.Error())
				}
			}

			if location.Pulau != "" && time.Tahun != "" {
				cube, err = helpers.CubeTimeLocation(time, location, other)
				if err != nil {
					return ResponseJson(ctx, http.StatusBadGateway, err.Error())
				}
			}
		} else {
			cube, err = helpers.SQLtoCube(other, limit, offset)
			if err != nil {
				return ResponseJson(ctx, http.StatusBadGateway, err.Error())
			}
		}
	} else {
		cube, err = helpers.SQLtoCube(other, limit, offset)
		if err != nil {
			return ResponseJson(ctx, http.StatusBadGateway, err.Error())
		}
	}

	if point != "" && dim != "" {
		cube = cube.Slice(dimension, point)
		cube = cube.RollUp([]string{dim}, cube.Fields, helpers.Sum, []interface{}{0})
	} else {
		cube = cube.RollUp([]string{dimension}, cube.Fields, helpers.Sum, []interface{}{0})
	}

	return ResponseJson(ctx, http.StatusOK, cube.Rows())
}

func handleLocationQueryWithDate(ctx *fiber.Ctx, selectedDate string) error {
    confidence := ctx.Query("confidence")
    satelite := ctx.Query("satelite")
    pulau := ctx.Query("pulau")
    provinsi := ctx.Query("provinsi")
    kota := ctx.Query("kota")
    kecamatan := ctx.Query("kecamatan")
    desa := ctx.Query("desa")

    fmt.Printf("handleLocationQueryWithDate - selectedDate: %s\n", selectedDate)
    fmt.Printf("Other filters - confidence: %s, satelite: %s, pulau: %s, provinsi: %s, kota: %s\n", 
        confidence, satelite, pulau, provinsi, kota)

    var locationData []struct {
        Location string `gorm:"column:location"`
        Total    int    `gorm:"column:total"`
    }

    // Menentukan level agregasi berdasarkan parameter drill-down
    var groupByColumn string = "dl.pulau" // Default ke pulau
    
    if desa != "" {
        groupByColumn = "dl.desa"
    } else if kecamatan != "" {
        groupByColumn = "dl.desa" // Group by desa dalam kecamatan
    } else if kota != "" {
        groupByColumn = "dl.kecamatan" // Group by kecamatan dalam kota
    } else if provinsi != "" {
        groupByColumn = "dl.kab_kota" // Group by kota dalam provinsi  
    } else if pulau != "" {
        groupByColumn = "dl.provinsi" // Group by provinsi dalam pulau
    }

    // Build dynamic query
    query := `
        SELECT 
            ` + groupByColumn + ` as location,
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
        query += ` AND DATE(fh.hotspot_time) = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, selectedDate)
        argIndex++
        fmt.Printf("Applied date filter to location query: %s\n", selectedDate)
    }

    // Apply other filters
    if confidence != "" {
        query += ` AND dc.confidence_level = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, confidence)
        argIndex++
    }
    
    if satelite != "" {
        query += ` AND ds.satelite_name = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, satelite)
        argIndex++
    }

    if pulau != "" {
        query += ` AND dl.pulau = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, pulau)
        argIndex++
    }

    if provinsi != "" {
        query += ` AND dl.provinsi = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, provinsi)
        argIndex++
    }

    if kota != "" {
        query += ` AND dl.kab_kota = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, kota)
        argIndex++
    }

    if kecamatan != "" {
        query += ` AND dl.kecamatan = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, kecamatan)
        argIndex++
    }

    if desa != "" {
        query += ` AND dl.desa = $` + fmt.Sprintf("%d", argIndex)
        args = append(args, desa)
        argIndex++
    }

    // Group by location
    query += ` GROUP BY ` + groupByColumn + ` ORDER BY total DESC`

    fmt.Printf("Final location query: %s\n", query)
    fmt.Printf("Query args: %v\n", args)

    // Execute query
    err := clients.DATABASE.Raw(query, args...).Scan(&locationData).Error
    if err != nil {
        fmt.Printf("Database error in handleLocationQueryWithDate: %v\n", err)
        return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
    }

    fmt.Printf("Location data retrieved with date filter: %d records\n", len(locationData))

    // Format response as [[location, total]]
    result := make([][]interface{}, 0)
    for _, item := range locationData {
        result = append(result, []interface{}{item.Location, item.Total})
    }

    fmt.Printf("Final result: %v\n", result)
    return ResponseJson(ctx, http.StatusOK, result)
}

func handleTimeQuery(ctx *fiber.Ctx) error {
    var timeParam payload.Time
    if err := ctx.QueryParser(&timeParam); err != nil {
        return ResponseJson(ctx, http.StatusBadRequest, "Invalid time parameters")
    }
    var results []string
    query := clients.DATABASE.Table("dim_time")

    if timeParam.Tahun != "" {
        query = query.Where("tahun = ?", timeParam.Tahun)
    }
    if timeParam.Semester != "" {
        query = query.Where("semester = ?", timeParam.Semester)
    }
    if timeParam.Kuartal != "" {
        query = query.Where("kuartal = ?", timeParam.Kuartal)
    }
    if timeParam.Bulan != "" {
        query = query.Where("bulan = ?", timeParam.Bulan)
    }
	if timeParam.Minggu != "" {
        query = query.Where("minggu = ?", timeParam.Minggu)
    }
    if timeParam.Hari != "" {
        query = query.Where("hari = ?", timeParam.Hari)
    }
    selectField := determineTimeField(timeParam)
    
    if err := query.Distinct(selectField).Pluck(selectField, &results).Error; err != nil {
        return ResponseJson(ctx, http.StatusInternalServerError, err.Error())
    }
    return ResponseJson(ctx, http.StatusOK, results)
}

func determineTimeField(param payload.Time) string {
    if param.Hari == "" && param.Minggu != "" {
        return "hari"
    }
	if param.Minggu == "" && param.Bulan != "" {
        return "minggu"
    }
    if param.Bulan == "" && param.Kuartal != "" {
        return "bulan"
    }
    if param.Kuartal == "" && param.Semester != "" {
        return "kuartal"
    }
    if param.Semester == "" && param.Tahun != "" {
        return "semester"
    }
    return "tahun" 
}