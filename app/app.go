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
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	configs.ConnectDB()
	routers.Router(app)
	log.Fatalln(app.Listen(os.Getenv("SERVER")))
}
