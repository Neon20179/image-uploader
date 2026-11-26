package pgx

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"image-uploader/internal/config"
)

type Config struct {
	Host     string        `envconfig:"HOST" default:"localhost"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"5s"`
	SSLMode  string        `envconfig:"SSL_MODE" default:"disable"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("POSTGRES", func(c Config) error {
		port, err := strconv.Atoi(c.Port)
		if strings.TrimSpace(c.Host) == "" || err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("host and port must form a valid PostgreSQL address")
		}
		if c.User == "" || c.Password == "" || c.Database == "" || c.Timeout <= 0 {
			return fmt.Errorf("credentials and database must not be empty, timeout must be positive")
		}
		switch c.SSLMode {
		case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
		default:
			return fmt.Errorf("unsupported SSL mode")
		}
		return nil
	})
}

func MustNewConfig() Config { return config.Must(NewConfig()) }
