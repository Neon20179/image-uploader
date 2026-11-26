package ocr

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	URL     string
	Timeout time.Duration
}

func (c Config) Validate() error {
	u, err := url.Parse(c.URL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return fmt.Errorf("OCR URL must be an HTTP(S) base URL")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("OCR timeout must be positive")
	}
	return nil
}
