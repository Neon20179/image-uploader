package server

import (
	"fmt"
	"time"

	"image-uploader/internal/config"
)

type Config struct {
	Addr              string        `envconfig:"ADDR" default:":8000"`
	ShutdownTimeout   time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"15s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"30s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("SERVER", func(c Config) error {
		if c.Addr == "" || c.ShutdownTimeout <= 0 || c.ReadHeaderTimeout <= 0 || c.ReadTimeout <= 0 || c.WriteTimeout <= 0 || c.IdleTimeout <= 0 {
			return fmt.Errorf("address must not be empty and timeouts must be positive")
		}
		return nil
	})
}

func MustNewConfig() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return config
}
