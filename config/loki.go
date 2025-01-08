package config

import (
	"fmt"

	"github.com/grafana/loki-client-go/loki"
)

func InitLoki(uri string) (*loki.Client, error) {
	if uri == "" {
		return nil, fmt.Errorf("empty env loki_url")
	}
	config, err := loki.NewDefaultConfig(uri)
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
