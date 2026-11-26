package minio

import (
	"fmt"
	"time"

	"image-uploader/internal/config"
)

type Config struct {
	Endpoint       string        `envconfig:"ENDPOINT" required:"true"`
	PublicEndpoint string        `envconfig:"PUBLIC_ENDPOINT"`
	PublicSecure   bool          `envconfig:"PUBLIC_SECURE" default:"false"`
	Region         string        `envconfig:"REGION" default:"us-east-1"`
	AccessKey      string        `envconfig:"ACCESS_KEY" required:"true"`
	SecretKey      string        `envconfig:"SECRET_KEY" required:"true"`
	Secure         bool          `envconfig:"SECURE" default:"false"`
	Timeout        time.Duration `envconfig:"TIMEOUT" default:"10s"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("MINIO", func(c Config) error {
		if c.Timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		return nil
	})
}

func MustNewConfig() Config { return config.Must(NewConfig()) }
