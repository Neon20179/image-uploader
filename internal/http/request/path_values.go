package request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
)

func GetStringPathValue(r *http.Request, key string) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", fmt.Errorf("no key=%s in path values: %w", key, apperr.ErrInvalidArgument)
	}
	return pathValue, nil
}

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue, err := GetStringPathValue(r, key)
	if err != nil {
		return uuid.Nil, err
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"path value=%s by key=%s not a valid UUID: %v: %w",
			pathValue,
			key,
			err,
			apperr.ErrInvalidArgument,
		)
	}

	return val, nil
}
