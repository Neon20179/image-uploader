package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"image-uploader/internal/apperr"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func Deserialize(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			apperr.ErrInvalidArgument,
		)
	}

	var err error

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			apperr.ErrInvalidArgument,
		)
	}

	return nil
}
