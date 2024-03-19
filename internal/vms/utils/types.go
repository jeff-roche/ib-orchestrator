package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNotImplemented = errors.New("functionality not yet implemented")
	SpecValidator     = validator.New(validator.WithRequiredStructEnabled())
)
