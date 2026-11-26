package finalizer

import (
	"image-uploader/internal/config"
)

type Config struct {
	Group          string `envconfig:"FINALIZER_GROUP" required:"true"`
	OCRTopic       string `envconfig:"OCR_TOPIC" required:"true"`
	EmbeddingTopic string `envconfig:"EMBEDDING_TOPIC" required:"true"`
	DLQTopic       string `envconfig:"FINALIZER_DLQ_TOPIC" required:"true"`
}

func NewConfig() (Config, error) {
	return config.Read[Config]("KAFKA")
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
