package ocr

import (
	"fmt"
	"math"
	"time"

	"image-uploader/internal/config"
	ocrclient "image-uploader/internal/ocr"
)

type ClientConfig struct {
	URL     string        `envconfig:"BASE_URL" required:"true"`
	Timeout time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func (cfg ClientConfig) Options() ocrclient.Config {
	return ocrclient.Config(cfg)
}

type Config struct {
	BlacklistWords               []string     `envconfig:"OCR_BLACKLIST_WORDS"`
	BlacklistSimilarityThreshold float64      `envconfig:"OCR_BLACKLIST_SIMILARITY_THRESHOLD" default:"0.80"`
	OCRTopic                     string       `envconfig:"KAFKA_OCR_TOPIC" required:"true"`
	OCR                          ClientConfig `envconfig:"OCR"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("", func(cfg Config) error { return cfg.Validate() })
}

func (cfg Config) Validate() error {
	threshold := cfg.BlacklistSimilarityThreshold
	if math.IsNaN(threshold) || math.IsInf(threshold, 0) || threshold <= 0 || threshold > 1 {
		return fmt.Errorf("OCR blacklist similarity threshold must be greater than 0 and at most 1")
	}
	for _, word := range cfg.BlacklistWords {
		if len(words(word)) != 1 {
			return fmt.Errorf("OCR blacklist entries must contain one word each")
		}
	}
	return cfg.OCR.Options().Validate()
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
