package config

import (
	_ "embed"
	"encoding/json"
	"github.com/pkg/errors"
)

type Database struct {
	DatabaseName string `json:"db_name"`
	User         string `json:"db_user"`
	Password     string `json:"db_pass"`
	URI          string `json:"db_dsn"`
}

type Config struct {
	Host            string    `json:"host"`
	Port            int       `json:"port"`
	ShutdownTimeout int       `json:"shutdown_timeout"`
	Database        *Database `json:"database"`
}

//go:embed tsconfig.json
var data []byte

// New returns a new configuration, and attempts to load
// config from file system.
func New() (*Config, error) {
	cfg := &Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {
		return nil, errors.Errorf("couldn't parse json file.: %s", err)
	}

	return cfg, nil
}
