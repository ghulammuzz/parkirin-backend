package config

import (
	"fmt"
	"os"

	"github.com/grafana/loki-client-go/loki"
)

func InitLoki() (*loki.Client, error) {
	lokiURL := os.Getenv("LOKI_URL")
	// slog.Info(lokiURL)
	if lokiURL == "" {
		lokiURL = "http://localhost:3100/loki/api/v1/push"
	}
	config, err := loki.NewDefaultConfig(lokiURL)
	if err != nil {
		return nil, fmt.Errorf("error in def config")
	}
	config.TenantID = "xyz"
	client, err := loki.New(config)
	if err != nil {
		return nil, fmt.Errorf("error in def new")
	}

	return client, nil
}
