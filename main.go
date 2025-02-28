package main

import (
	"fib/database"
	"fib/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectDB()

	app := fiber.New(
		fiber.Config{
			Prefork:       true,
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Fiber",
			AppName:       "App name",
		},
	)

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
