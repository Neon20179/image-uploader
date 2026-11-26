package image

import (
	"fmt"
	"net/url"

	"image-uploader/internal/config"
)

type Config struct {
	PublicBaseURL string `envconfig:"MINIO_BASE_URL" required:"true"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("", func(c Config) error {
		u, err := url.Parse(c.PublicBaseURL)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return fmt.Errorf("public image base URL must be an absolute HTTP(S) URL")
		}
		return nil
	})
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}
