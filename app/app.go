package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/nechitast/olap-backend/app/configs"
	"github.com/nechitast/olap-backend/app/routers"
)

func RunApp() {
	godotenv.Load()
	app := fiber.New()

	// CORS configuration fixed for production
	// Only allow frontend origin to prevent CORS errors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://app.olaphotspot.web.id",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	configs.ConnectDB()
	routers.Router(app)
	log.Fatalln(app.Listen(os.Getenv("SERVER")))
}
