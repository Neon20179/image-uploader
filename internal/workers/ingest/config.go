package ingest

import (
	"image-uploader/internal/config"
)

type Config struct {
	Group    string `envconfig:"INGEST_GROUP" required:"true"`
	Topic    string `envconfig:"MINIO_TOPIC" required:"true"`
	DLQTopic string `envconfig:"INGEST_DLQ_TOPIC" required:"true"`
}

func NewConfig() (Config, error) {
	return config.Read[Config]("KAFKA")
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
