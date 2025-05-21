package configs

import (
	"log"
	"os"

	"github.com/nechitast/olap-backend/app/configs/clients"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error
	// https: //github.com/go-gorm/postgres
	clients.DATABASE, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DATABASE_DSN"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	// clients.DATABASE, err = gorm.Open(mysql.New(mysql.Config{
	// 	DSN:                       os.Getenv("DATABASE_DSN"), // data source name
	// 	DefaultStringSize:         256,                       // default size for string fields
	// 	DisableDatetimePrecision:  true,                      // disable datetime precision, which not supported before MySQL 5.6
	// 	DontSupportRenameIndex:    true,                      // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
	// 	DontSupportRenameColumn:   true,                      // `change` when rename column, rename column not supported before MySQL 8, MariaDB
	// 	SkipInitializeWithVersion: false,                     // auto configure based on currently MySQL version
	// }), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// clients.DATABASE.AutoMigrate(
	// 	models.Dim_Confidence{},
	// 	models.Fact_Hotspot{},
	// 	models.Dim_Location{},
	// 	models.Dim_Satelite{},
	// 	models.Dim_Time{},
	// )

}
