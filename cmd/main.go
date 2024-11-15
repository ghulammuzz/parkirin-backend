package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ghulammuzz/backend-parkerin/config"
	applicants "github.com/ghulammuzz/backend-parkerin/internal/applicants/di"
	"github.com/ghulammuzz/backend-parkerin/internal/health"
	store "github.com/ghulammuzz/backend-parkerin/internal/store/di"
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

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	// app.Use(recover.New())

	app.Get("/hc", health.HealthCheck(db))

	api := app.Group("/api")
	users.InitializedUsersService(db, config.Validate).Router(api)
	store.InitializedStoreService(db).Router(api)
	applicants.InitializedApplicationService(db).Router(api)

	if err := app.Listen(fmt.Sprint(":", os.Getenv("APP_PORT"))); err != nil {
		log.Error("Failed to start the server: %v", err)
	}

}
