package helpers

import (
	"log"
	"github.com/nechitast/olap-backend/app/configs/clients"
	"github.com/nechitast/olap-backend/app/models"
	"github.com/nechitast/olap-backend/app/models/payload"
	"github.com/xdbsoft/olap"
)

func SQLtoCube(oth payload.Other, limit int, offset int) (olap.Cube, error) {
	var count int64
	var list []models.Fact_Hotspot
	var data [][]interface{}
	// var dataLocation [][]interface{}

	cube := olap.Cube{
		Dimensions: []string{"location", "time", "confidence", "satelite"},
		Fields:     []string{"total"},
	}

	where := make(map[string]interface{})

	if oth.Confidence != "" {
		log.Println("hello")
		where["Dim_Confidence.confidence_level"] = oth.Confidence
	}

	if oth.Satelite != "" {
		where["Dim_Satelite.satelite_name"] = oth.Satelite
	}

	if err := clients.DATABASE.Table("fact_hotspot").
		Preload("Dim_Location").Joins("Dim_Confidence").
		Preload("Dim_Time").Joins("Dim_Satelite").Where(where).
		Find(&list).Count(&count).Error; err != nil {
		return cube, err
	}

	for i := 0; i < int(count); i++ {
		item := list[i]
		data = append(data, []interface{}{
			item.Dim_Location.Pulau,
			item.Dim_Time.Tahun,
			item.Dim_Confidence.Confidence_Level,
			item.Dim_Satelite.Satelite_Name,
			item.Hotspot_Count,
		})
	}

	cube.AddRows([]string{"location", "time", "confidence", "satelite", "total"}, data)

	log.Println(cube.Headers())
	return cube, nil
}

func CubeLocation(param payload.Location, oth payload.Other) (olap.Cube, error) {
	var count int64
	var list []models.Fact_Hotspot
	var data [][]interface{}

	cube := olap.Cube{
		Dimensions: []string{"location", "time", "confidence", "satelite"},
		Fields:     []string{"total"},
	}

	where := make(map[string]interface{})

	if oth.Confidence != "" {
		where["Dim_Confidence.confidence_level"] = oth.Confidence
	}

	if oth.Satelite != "" {
		where["Dim_Satelite.satelite_name"] = oth.Satelite
	}

	if param.Kecamatan != "" {
		where["Dim_Location.kecamatan"] = param.Kecamatan
	}

	if param.Kota != "" {
		where["Dim_Location.kab_kota"] = param.Kota
	}

	if param.Provinsi != "" {
		where["Dim_Location.provinsi"] = param.Provinsi
	}

	if param.Pulau != "" {
		where["Dim_Location.pulau"] = param.Pulau
	}

	print(where["Dim_Location.kab_kota"])

	if err := clients.DATABASE.Table("fact_hotspot").
		Joins("Dim_Location").Joins("Dim_Confidence").
		Preload("Dim_Time").Joins("Dim_Satelite").Where(where).
		Find(&list).Count(&count).Error; err != nil {
		return cube, err
	}

	for i := 0; i < int(count); i++ {
		var value string
		item := list[i]
		if param.Kecamatan != "" {
			value = item.Dim_Location.Desa
		} else if param.Kota != "" {
			value = item.Dim_Location.Kecamatan
		} else if param.Provinsi != "" {
			value = item.Dim_Location.Kab_kota
		} else if param.Pulau != "" {
			value = item.Dim_Location.Provinsi
		} else {
			value = item.Dim_Location.Pulau
		}
		data = append(data, []interface{}{
			value,
			item.Dim_Time.Tahun,
			item.Dim_Confidence.Confidence_Level,
			item.Dim_Satelite.Satelite_Name,
			item.Hotspot_Count,
		})
	}

	cube.AddRows([]string{"location", "time", "confidence", "satelite", "total"}, data)

	return cube, nil
}

func CubeTime(param payload.Time, oth payload.Other) (olap.Cube, error) {
	var count int64
	var list []models.Fact_Hotspot
	var data [][]interface{}

	cube := olap.Cube{
		Dimensions: []string{"location", "time", "confidence", "satelite"},
		Fields:     []string{"total"},
	}

	where := make(map[string]interface{})

	if oth.Confidence != "" {
		where["Dim_Confidence.confidence_level"] = oth.Confidence
	}

	if oth.Satelite != "" {
		where["Dim_Satelite.satelite_name"] = oth.Satelite
	}

	if param.Hari != "" {
		where["Dim_Time.hari"] = param.Hari
	}

	if param.Minggu != "" {
		where["Dim_Time.minggu"] = param.Minggu
	}

	if param.Bulan != "" {
		where["Dim_Time.bulan"] = param.Bulan
	}

	if param.Kuartal != "" {
		where["Dim_Time.kuartal"] = param.Kuartal
	}

	if param.Semester != "" {
		where["Dim_Time.semester"] = param.Semester
	}

	if param.Tahun != "" {
		where["Dim_Time.tahun"] = param.Tahun
	}

	if err := clients.DATABASE.Table("fact_hotspot").
		Preload("Dim_Location").Joins("Dim_Confidence").
		Joins("Dim_Time").Joins("Dim_Satelite").Where(where).
		Limit(limit).Offset(offset).Find(&list).Count(&count).Error; err != nil {
		return cube, err
	}

	for i := 0; i < int(count); i++ {
		var value string
		item := list[i]
		if param.Minggu != "" {
			value = item.Dim_Time.Hari
		} else if param.Bulan != "" {
			value = (string)(item.Dim_Time.Minggu)
		} else if param.Kuartal != "" {
			value = item.Dim_Time.Bulan
		} else if param.Semester != "" {
			value = item.Dim_Time.Kuartal
		} else if param.Tahun != "" {
			value = (string)(item.Dim_Time.Semester)
		}
		data = append(data, []interface{}{
			item.Dim_Location.Pulau,
			value,
			item.Dim_Confidence.Confidence_Level,
			item.Dim_Satelite.Satelite_Name,
			item.Hotspot_Count,
		})
	}

	cube.AddRows([]string{"location", "time", "confidence", "satelite", "total"}, data)

	return cube, nil
}

func CubeTimeLocation(param payload.Time, loc payload.Location, oth payload.Other) (olap.Cube, error) {
	var count int64
	var list []models.Fact_Hotspot
	var data [][]interface{}

	cube := olap.Cube{
		Dimensions: []string{"location", "time", "confidence", "satelite"},
		Fields:     []string{"total"},
	}

	where := make(map[string]interface{})
	whereLoc := make(map[string]interface{})

	if oth.Confidence != "" {
		where["Dim_Confidence.confidence_level"] = oth.Confidence
	}

	if oth.Satelite != "" {
		where["Dim_Satelite.satelite_name"] = oth.Satelite
	}

	if param.Hari != "" {
		where["Dim_Time.hari"] = param.Hari
	}

	if param.Minggu != "" {
		where["Dim_Time.minggu"] = param.Minggu
	}

	if param.Bulan != "" {
		where["Dim_Time.bulan"] = param.Bulan
	}

	if param.Kuartal != "" {
		where["Dim_Time.kuartal"] = param.Kuartal
	}

	if param.Semester != "" {
		where["Dim_Time.semester"] = param.Semester
	}

	if param.Tahun != "" {
		where["Dim_Time.tahun"] = param.Tahun
	}

	if loc.Kecamatan != "" {
		whereLoc["Dim_Location.kecamatan"] = loc.Kecamatan
	}

	if loc.Kota != "" {
		whereLoc["Dim_Location.kab_kota"] = loc.Kota
	}

	if loc.Provinsi != "" {
		whereLoc["Dim_Location.provinsi"] = loc.Provinsi
	}

	if loc.Pulau != "" {
		whereLoc["Dim_Location.pulau"] = loc.Pulau
	}

	if err := clients.DATABASE.Table("fact_hotspot").
		Joins("Dim_Location").Joins("Dim_Confidence").
		Joins("Dim_Time").Joins("Dim_Satelite").Where(where).Where(whereLoc).
		Find(&list).Count(&count).Error; err != nil {
		return cube, err
	}

	for i := 0; i < int(count); i++ {
		var time string
		var location string

		item := list[i]

		if param.Minggu != "" {
			time = item.Dim_Time.Hari
		} else if param.Bulan != "" {
			time = (string)(item.Dim_Time.Minggu)
		} else if param.Kuartal != "" {
			time = item.Dim_Time.Bulan
		} else if param.Semester != "" {
			time = item.Dim_Time.Kuartal
		} else if param.Tahun != "" {
			time = (string)(item.Dim_Time.Semester)
		} else {
			time = (string)(item.Dim_Time.Tahun)
		}

		if loc.Kecamatan != "" {
			location = item.Dim_Location.Desa
		} else if loc.Kota != "" {
			location = item.Dim_Location.Kecamatan
		} else if loc.Provinsi != "" {
			location = item.Dim_Location.Kab_kota
		} else if loc.Pulau != "" {
			location = item.Dim_Location.Provinsi
		} else {
			location = item.Dim_Location.Pulau
		}
		data = append(data, []interface{}{
			location,
			time,
			item.Dim_Confidence.Confidence_Level,
			item.Dim_Satelite.Satelite_Name,
			item.Hotspot_Count,
		})
	}

	cube.AddRows([]string{"location", "time", "confidence", "satelite", "total"}, data)

	return cube, nil
}

func Sum(aggregate, value []interface{}) []interface{} {
	s := aggregate[0].(int)
	s += value[0].(int)
	return []interface{}{s}
}
