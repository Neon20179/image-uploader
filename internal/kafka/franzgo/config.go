package franzgo

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"image-uploader/internal/config"
)

type Config struct {
	Host            string        `envconfig:"HOST" default:"localhost"`
	Port            string        `envconfig:"PORT" default:"9092"`
	MaxAttempts     int           `envconfig:"MAX_ATTEMPTS" default:"3"`
	RetryTimeout    time.Duration `envconfig:"RETRY_TIMEOUT" default:"1s"`
	WaitTimeout     time.Duration `envconfig:"WAIT_TIMEOUT" default:"5s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
	DeliveryTimeout time.Duration `envconfig:"DELIVERY_TIMEOUT" default:"5s"`
	Linger          time.Duration `envconfig:"LINGER" default:"10ms"`
}

func NewConfig() (Config, error) {
	return config.ReadAndValidate("KAFKA", func(c Config) error {
		port, err := strconv.Atoi(c.Port)
		if strings.TrimSpace(c.Host) == "" || err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("host and port must form a valid broker address")
		}
		if c.MaxAttempts <= 0 || c.RetryTimeout <= 0 || c.WaitTimeout <= 0 || c.ShutdownTimeout <= 0 || c.DeliveryTimeout <= 0 || c.Linger < 0 {
			return fmt.Errorf("attempts and timeouts must be positive; linger must be nonnegative")
		}
		return nil
	})
}

func MustNewConfig() Config {
	return config.Must(NewConfig())
}

func (c Config) Broker() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c Config) Consumer(group string, topics []string, dlqTopic string) ConsumerConfig {
	return ConsumerConfig{
		Broker:          c.Broker(),
		Group:           group,
		Topics:          topics,
		DLQTopic:        dlqTopic,
		MaxAttempts:     c.MaxAttempts,
		RetryTimeout:    c.RetryTimeout,
		WaitTimeout:     c.WaitTimeout,
		DeliveryTimeout: c.DeliveryTimeout,
		ShutdownTimeout: c.ShutdownTimeout,
	}
}

func (c Config) Producer() ProducerConfig {
	return ProducerConfig{
		Broker:          c.Broker(),
		DeliveryTimeout: c.DeliveryTimeout,
		Linger:          c.Linger,
		ShutdownTimeout: c.ShutdownTimeout,
	}
}
