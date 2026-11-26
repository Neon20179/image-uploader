package logger

import (
	"fmt"
	"log/slog"

	"image-uploader/internal/config"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"info"`
	Folder string `envconfig:"FOLDER" default:"out/logs"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("LOGGER", func(cfg Config) error {
		var level slog.Level
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			return fmt.Errorf("invalid log level")
		}
		return nil
	})
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
