package main

import (
	"flag"

	"github.com/ghulammuzz/backend-parkerin/config"
	"github.com/ghulammuzz/backend-parkerin/internal/health"
	users "github.com/ghulammuzz/backend-parkerin/internal/users/di"
	"github.com/ghulammuzz/backend-parkerin/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	env := flag.String("env", "prod", "Environment for (stg/prod)")
	flag.Parse()

	if *env == "stg" {
		log.InitLogger("dev")
		err := godotenv.Load("./stg.env")
		if err != nil {
			log.Error("Error loading stg.env file: %v", err)
		}
		log.Info("Environment: staging (stg.env loaded)")
	} else {
		log.InitLogger("prod")
		log.Info("Environment: production (using system environment variables)")
	}
	config.InitStorage()
	config.InitValidator()
}

func main() {

	db, err := config.InitPostgres()
	if err != nil {
		log.Error("Failed to initialize database: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	app.Get("/hc", health.HealthCheck(db))

	api := app.Group("/api")
	users.InitializedUsersService(db, config.Validate).Router(api)

	if err := app.Listen(":3000"); err != nil {
		log.Error("Failed to start the server: %v", err)
	}

}
