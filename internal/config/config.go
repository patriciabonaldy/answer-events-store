package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Host            string        `env:"HOST"`
	Port            int           `env:"PORT"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIME"`
	MongoConfig     *MongoConfig
}

type MongoConfig struct {
	URI          string `env:"DATABASE_DSN"`
	DatabaseName string `env:"DATABASE_NAME"`
	User         string `env:"DATABASE_USR"`
	Password     string `env:"DATABASE_PASS"`
}

// New returns a new configuration, and attempts to load
// config from file system.
func New() (*Config, error) {
	env := strings.ToLower(os.Getenv("GO_ENVIRONMENT"))
	viper.SetConfigFile(fmt.Sprintf("%s.env", env))
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var vipErr viper.ConfigFileNotFoundError
		if ok := errors.As(err, &vipErr); ok {
			log.Fatalln(fmt.Errorf("config file not found. %w", err))
		} else {
			log.Fatalln(fmt.Errorf("unexpected error loading config file. %w", err))
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln(fmt.Errorf("failed to unmarshal config. %w", err))
	}

	return config, nil
}
