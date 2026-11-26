package embedding

import (
	"image-uploader/internal/config"
)

type Config struct {
	Group    string `envconfig:"EMBEDDING_GROUP" required:"true"`
	Topic    string `envconfig:"NORMALIZED_TOPIC" required:"true"`
	DLQTopic string `envconfig:"EMBEDDING_DLQ_TOPIC" required:"true"`
}

func NewConfig() (Config, error) {
	return config.Read[Config]("KAFKA")
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
