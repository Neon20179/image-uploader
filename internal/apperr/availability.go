package apperr

import "errors"

type UnavailableError struct {
	Err error
}

func (e *UnavailableError) Error() string {
	return e.Err.Error()
}

func (e *UnavailableError) Unwrap() error {
	return e.Err
}

func NewUnavailableError(err error) error {
	if err == nil {
		return nil
	}
	return &UnavailableError{Err: err}
}

func IsUnavailableError(err error) bool {
	var unavailable *UnavailableError
	return errors.As(err, &unavailable)
}
