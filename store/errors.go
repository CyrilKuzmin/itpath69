// Store defines the main interface for the storage

package store

import (
	"errors"
	"fmt"
)

// ErrorType define the kind of service's errors
type ErrorType int8

// Service error
const (
	InternalErr ErrorType = iota
	NotFoundErr
	AlreadyExistsErr
	InvalidArgsErr
	ForbiddenErr
	UnavailableErr
)

type StoreError struct {
	etype ErrorType
	msg   string
	error error
}

func (e StoreError) Error() string {
	return e.msg
}

// Unwrap unwraps an error
func (e StoreError) Unwrap() error {
	return e.error
}

// ErrorIs checks if given error has specified ErrorType.
// It returns false if error is not service Error at all.
func ErrorIs(err error, etype ErrorType) bool {
	var storeErr *StoreError
	if !errors.As(err, &storeErr) {
		return false
	}
	return storeErr.etype == etype
}

func newError(etype ErrorType, wrappedErr error, msg string) error {
	if wrappedErr != nil {
		if msg != "" {
			msg += ": "
		}
		msg += wrappedErr.Error()
	}
	return &StoreError{
		etype: etype,
		msg:   msg,
		error: wrappedErr,
	}
}

// NewError returns new service Error with type and message specified
var NewError = newError

func newErrorf(etype ErrorType, wrappedErr error, format string, args ...interface{}) error {
	return newError(etype, wrappedErr, fmt.Sprintf(format, args...))
}

func ErrInternal(err error) error {
	return newErrorf(InternalErr, err, "internal storage error")
}

func ErrUserAlreadyExists(username string) error {
	return newErrorf(AlreadyExistsErr, nil, "user %v already exists", username)
}

func ErrUserNotFound(username string) error {
	return newErrorf(NotFoundErr, nil, "user %v not found", username)
}

func ErrModuleNotFound(id int) error {
	return newErrorf(NotFoundErr, nil, "module %v not found", id)
}
