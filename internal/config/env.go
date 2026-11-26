package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func Read[T any](prefix string) (T, error) {
	return ReadAndValidate(prefix, func(t T) error { return nil })
}

func ReadAndValidate[T any](prefix string, validate func(T) error) (T, error) {
	var cfg T
	if err := envconfig.Process(prefix, &cfg); err != nil {
		if parseErr, ok := err.(*envconfig.ParseError); ok {
			return cfg, fmt.Errorf("%s: invalid %s (expected %s)", prefix, parseErr.KeyName, parseErr.TypeName)
		}
		return cfg, fmt.Errorf("%s config: %w", prefix, err)
	}
	if err := validate(cfg); err != nil {
		return cfg, fmt.Errorf("%s config: %w", prefix, err)
	}
	return cfg, nil
}

func Must[T any](cfg T, err error) T {
	if err != nil {
		panic(err)
	}
	return cfg
}
