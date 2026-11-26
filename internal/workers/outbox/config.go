package outbox

import (
	"fmt"
	"time"

	"image-uploader/internal/config"
)

type Config struct {
	BatchSize    int           `envconfig:"BATCH_SIZE" default:"100"`
	PollInterval time.Duration `envconfig:"POLL_INTERVAL" default:"1s"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("OUTBOX", func(cfg Config) error {
		if cfg.BatchSize < 0 {
			return fmt.Errorf("OUTBOX_BATCH_SIZE must be positive")
		}
		if cfg.PollInterval < 0 {
			return fmt.Errorf("OUTBOX_POLL_INTERVAL must be positive")
		}
		return nil
	})
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
