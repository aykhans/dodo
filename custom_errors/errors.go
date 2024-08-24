package customerrors

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidJSON = errors.New("invalid JSON file")
	ErrInvalidFile = errors.New("invalid file")
	ErrInterrupt   = errors.New("interrupted")
	ErrNoInternet  = errors.New("no internet connection")
)

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

type Error interface {
	Error() string
	Unwrap() error
}

type TypeError struct {
	Expected string
	Received string
	Field    string
	err      error
}

func NewTypeError(expected, received, field string, err error) *TypeError {
	return &TypeError{
		Expected: expected,
		Received: received,
		Field:    field,
		err:      err,
	}
}

func (e *TypeError) Error() string {
	return "Expected " + e.Expected + " but received " + e.Received + " in field " + e.Field
}

func (e *TypeError) Unwrap() error {
	return e.err
}

type InvalidFileError struct {
	FileName string
	err      error
}

func NewInvalidFileError(fileName string, err error) *InvalidFileError {
	return &InvalidFileError{
		FileName: fileName,
		err:      err,
	}
}

func (e *InvalidFileError) Error() string {
	return "Invalid file: " + e.FileName
}

func (e *InvalidFileError) Unwrap() error {
	return e.err
}

type FileNotFoundError struct {
	FileName string
	err      error
}

func NewFileNotFoundError(fileName string, err error) *FileNotFoundError {
	return &FileNotFoundError{
		FileName: fileName,
		err:      err,
	}
}

func (e *FileNotFoundError) Error() string {
	return "File not found: " + e.FileName
}

func (e *FileNotFoundError) Unwrap() error {
	return e.err
}

type ValidationErrors struct {
	MapErrors map[string]string
	errors    validator.ValidationErrors
}

func NewValidationErrors(errsMap map[string]string, errs validator.ValidationErrors) *ValidationErrors {
	return &ValidationErrors{
		MapErrors: errsMap,
		errors:    errs,
	}
}

func (errs *ValidationErrors) Error() string {
	var errorsStr string
	for k, v := range errs.MapErrors {
		errorsStr += fmt.Sprintf("[%s]: %s\n", k, v)
	}
	return errorsStr
}

func (errs *ValidationErrors) Unwrap() error {
	return errs.errors
}
