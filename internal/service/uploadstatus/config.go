package uploadstatus

import (
	"image-uploader/internal/config"
)

type Config struct {
	UploadedTopic string `envconfig:"KAFKA_UPLOADED_TOPIC" required:"true"`
}

func NewConfig() (Config, error) {
	return config.Read[Config]("")
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
