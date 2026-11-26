package upload

import (
	"fmt"
	"time"

	"image-uploader/internal/config"
	"image-uploader/internal/imageproc"
)

type ImageConfig struct {
	MaxWidth       int `envconfig:"MAX_WIDTH" default:"1920"`
	MaxHeight      int `envconfig:"MAX_HEIGHT" default:"1920"`
	MaxSizeBytes   int `envconfig:"MAX_SIZE_BYTES" default:"1048576"`
	InitialQuality int `envconfig:"INITIAL_QUALITY" default:"85"`
	MinQuality     int `envconfig:"MIN_QUALITY" default:"40"`
	QualityStep    int `envconfig:"QUALITY_STEP" default:"5"`
	EncodeMethod   int `envconfig:"ENCODE_METHOD" default:"4"`
}

func (cfg ImageConfig) Options() imageproc.Config {
	return imageproc.Config(cfg)
}

type Config struct {
	Bucket                    string        `envconfig:"MINIO_BUCKET" required:"true"`
	RawObjectKeyPrefix        string        `envconfig:"UPLOAD_RAW_OBJECT_KEY_PREFIX" required:"true"`
	NormalizedObjectKeyPrefix string        `envconfig:"UPLOAD_NORMALIZED_OBJECT_KEY_PREFIX" required:"true"`
	NormalizedTopic           string        `envconfig:"KAFKA_NORMALIZED_TOPIC" required:"true"`
	MaxRawSizeBytes           int64         `envconfig:"MINIO_MAX_UPLOAD_SIZE_BYTES" default:"20971520"`
	MinRawSizeBytes           int64         `envconfig:"MINIO_MIN_UPLOAD_SIZE_BYTES" default:"1024"`
	URLExpiry                 time.Duration `envconfig:"UPLOAD_URL_EXPIRY" default:"1h"`
	NormalizeTimeout          time.Duration `envconfig:"UPLOAD_NORMALIZE_TIMEOUT" default:"10s"`
	AllowedMIMETypes          []string      `envconfig:"UPLOAD_ALLOWED_MIME_TYPES" default:"image/jpeg,image/png"`
	Image                     ImageConfig   `envconfig:"IMAGE"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("", func(cfg Config) error {
		if cfg.MinRawSizeBytes <= 0 {
			return fmt.Errorf("minimum raw image size must be positive")
		}
		if cfg.MaxRawSizeBytes < cfg.MinRawSizeBytes {
			return fmt.Errorf("maximum must be at least minimum")
		}
		return cfg.Image.Options().Validate()
	})
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
