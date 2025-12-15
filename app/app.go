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
		AllowOriginsFunc: func(origin string) bool {
			return origin == "https://app.olaphotspot.web.id" ||
				origin == "http://localhost:3000"
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	configs.ConnectDB()
	routers.Router(app)
	log.Fatalln(app.Listen(os.Getenv("SERVER")))
}
