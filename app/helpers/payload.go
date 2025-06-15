package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nechitast/olap-backend/app/models/payload"
)

func GetPayload(ctx *fiber.Ctx) (payload.Location, payload.Time, payload.Other) {
	var location payload.Location
	var time payload.Time
	var other payload.Other

	location.Pulau = ctx.Query("pulau", "")
	location.Provinsi = ctx.Query("provinsi", "")
	location.Kota = ctx.Query("kota", "")
	location.Kecamatan = ctx.Query("kecamatan", "")

	time.Tahun = ctx.Query("tahun", "")
	time.Semester = ctx.Query("semester", "")
	time.Kuartal = ctx.Query("kuartal", "")
	time.Bulan = ctx.Query("bulan", "")
	time.Minggu = ctx.Query("minggu", "")
	time.Hari = ctx.Query("hari", "")

	other.Confidence = ctx.Query("confidence", "")
	other.Satelite = ctx.Query("satelite", "")

	return location, time, other
}
