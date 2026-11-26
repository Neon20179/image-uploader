package middleware

import (
	"fmt"
	"net/url"

	"image-uploader/internal/config"
)

type Config struct {
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:3000"`
	AllowedMethods []string `envconfig:"ALLOWED_METHODS" default:"GET,POST,PATCH,DELETE,OPTIONS"`
	AllowedHeaders []string `envconfig:"ALLOWED_HEADERS" default:"Content-Type,Authorization"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("CORS", func(c Config) error {
		for _, origin := range c.AllowedOrigins {
			u, err := url.Parse(origin)
			if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
				return fmt.Errorf("allowed origins must be HTTP(S) origins without paths or credentials")
			}
		}
		return nil
	})
}

func MustNewConfig() Config { return config.Must(NewConfig()) }

type RequestLimitsConfig struct {
	RequestsPerSecond int `envconfig:"RATE_LIMIT_RPS" default:"40"`
	Burst             int `envconfig:"RATE_LIMIT_BURST" default:"60"`
	MaxInFlight       int `envconfig:"MAX_IN_FLIGHT" default:"20"`
}

func MustNewRequestLimitsConfig() RequestLimitsConfig {
	return config.Must(config.ReadAndValidate("API", func(c RequestLimitsConfig) error {
		if c.RequestsPerSecond <= 0 || c.Burst <= 0 || c.MaxInFlight <= 0 {
			return fmt.Errorf("RATE_LIMIT_RPS, RATE_LIMIT_BURST, and MAX_IN_FLIGHT must be positive")
		}
		return nil
	}))
}
