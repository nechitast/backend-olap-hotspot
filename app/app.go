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

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://app.olaphotspot.web.id,http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Content-Type",
	}))
	configs.ConnectDB()
	routers.Router(app)
	log.Fatalln(app.Listen(os.Getenv("SERVER")))
}
