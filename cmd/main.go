package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ghulammuzz/backend-parkerin/config"
	applicants "github.com/ghulammuzz/backend-parkerin/internal/applicants/di"
	health "github.com/ghulammuzz/backend-parkerin/internal/health"
	store "github.com/ghulammuzz/backend-parkerin/internal/store/di"
	users "github.com/ghulammuzz/backend-parkerin/internal/users/di"
	"github.com/gofiber/fiber/v2"

	mlog "log/slog"

	"github.com/ghulammuzz/backend-parkerin/pkg/log"
	"github.com/joho/godotenv"
)

func init() {
	env := flag.String("env", "prod", "Environment for (stg/prod)")
	flag.Parse()

	if *env == "stg" {
		err := godotenv.Load("./stg.env")
		if err != nil {
			mlog.Error("Error loading stg.env file ")
		}
		mlog.Info("Environment: staging (stg.env loaded)")
		mlog.Debug("debug tests")
	} else {
		mlog.Info("Environment: production (using system environment variables)")
	}

	lokiClient, err := config.InitLoki()
	if err != nil {
		fmt.Println("Error init loki")
		return
	}

	if *env == "stg" {
		log.InitLogger("dev", lokiClient)
	} else {
		log.InitLogger("prod", lokiClient)
	}

	config.InitStorage()
	config.InitValidator()
}

func main() {
	db, err := config.InitPostgres()
	if err != nil {
		log.Error("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// midtransClient := config.InitMidtrans()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/hc", health.HealthCheck(db))

	api := app.Group("/api")
	users.InitializedUsersService(db, config.Validate).Router(api)
	store.InitializedStoreService(db).Router(api)
	applicants.InitializedApplicationService(db).Router(api)
	// payment.InitializedPaymentService(db, midtransClient).Router(api)

	if err := app.Listen(fmt.Sprint(":", os.Getenv("APP_PORT"))); err != nil {
		log.Error("Failed to start the server: %v", err)
	}
}
