package kafka

import "errors"

type PermanentError struct {
	Err error
}

func (e *PermanentError) Error() string {
	return e.Err.Error()
}

func (e *PermanentError) Unwrap() error {
	return e.Err
}

func NewPermanentError(err error) error {
	if err == nil {
		return nil
	}
	return &PermanentError{Err: err}
}

func IsPermanentError(err error) bool {
	var pe *PermanentError
	return errors.As(err, &pe)
}

type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

func NewRetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &RetryableError{Err: err}
}

func IsRetryableError(err error) bool {
	var retryable *RetryableError
	return errors.As(err, &retryable)
}

type WaitError struct {
	Err error
}

func (e *WaitError) Error() string {
	return e.Err.Error()
}

func (e *WaitError) Unwrap() error {
	return e.Err
}

func NewWaitError(err error) error {
	if err == nil {
		return nil
	}
	return &WaitError{Err: err}
}

func IsWaitError(err error) bool {
	var wait *WaitError
	return errors.As(err, &wait)
}
