package embedding

import (
	"fmt"
	"math"
	"time"

	"image-uploader/internal/clip"
	"image-uploader/internal/config"
)

type ClipConfig struct {
	URL     string        `envconfig:"BASE_URL" required:"true"`
	Timeout time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func (cfg ClipConfig) Options() clip.Config {
	return clip.Config(cfg)
}

type Config struct {
	BlacklistSimilarityThreshold float64    `envconfig:"EMBEDDING_BLACKLIST_SIMILARITY_THRESHOLD" default:"0.90"`
	EmbeddingTopic               string     `envconfig:"KAFKA_EMBEDDING_TOPIC" required:"true"`
	Clip                         ClipConfig `envconfig:"CLIP"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("", func(cfg Config) error {
		return cfg.Validate()
	})
}

func (cfg Config) Validate() error {
	threshold := cfg.BlacklistSimilarityThreshold
	if math.IsNaN(threshold) || math.IsInf(threshold, 0) || threshold < -1 || threshold > 1 {
		return fmt.Errorf("embedding blacklist similarity threshold must be between -1 and 1")
	}
	return nil
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
