package api

import (
	v "github.com/go-playground/validator/v10"
)

// Validator represents a struct that wraps the go-playground validator.
// It is used to validate structs based on validation tags.
type Validator struct {
	Validator *v.Validate
}

// Validate validates the given interface 'i' using the configured validator.
// It returns an error if the validation fails or nil if the validation succeeds.
func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
