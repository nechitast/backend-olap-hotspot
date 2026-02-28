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
		// AllowOrigins:     "https://app.olaphotspot.web.id,http://localhost:3000",
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	configs.ConnectDB()
	routers.Router(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
	// log.Fatalln(app.Listen(os.Getenv("SERVER")))
}
