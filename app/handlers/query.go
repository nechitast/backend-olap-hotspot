package handlers

import (
	"net/http"
	"github.com/nechitast/olap-backend/app/configs/clients"
	"github.com/gofiber/fiber/v2"
	"github.com/nechitast/olap-backend/app/helpers"
	"github.com/nechitast/olap-backend/app/models/payload"
	"github.com/xdbsoft/olap"
	
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
    if param.Hari == "" && param.Bulan != "" {
        return "hari"
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