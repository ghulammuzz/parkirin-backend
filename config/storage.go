package config

import (
	"context"
	"encoding/base64"
	"os"

	"cloud.google.com/go/storage"
	"github.com/ghulammuzz/backend-parkerin/pkg/log"
	"google.golang.org/api/option"
)

var GCSClient *storage.Client

func InitStorage() {
	ctx := context.Background()

	credJSON, err := base64.StdEncoding.DecodeString(os.Getenv("GCS_CREDENTIALS_BASE64"))
	if err != nil {
		log.Error("error decoding credentials", err.Error())
		return
	}

	// log.Info(string(credJSON))

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credJSON))
	if err != nil {
		log.Error("error initializing storage", err.Error())
		return
	}

	GCSClient = client

}
