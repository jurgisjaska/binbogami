package api

import (
	v "github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *v.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
