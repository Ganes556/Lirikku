package middlewares

import (
	"errors"
	"strings"

	"github.com/go-playground/validator"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{Validator: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {	
	if err := v.Validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			if err.Tag() == "required" {
				if field == "artistnames"{
					return errors.New("artist_names is required")
				}
				return errors.New(field + " is required")
			}
			if err.Tag() == "email" {
				return errors.New(field + " is not valid")
			}
			if err.Tag() == "min" {
				return errors.New(field + " must be at least " + err.Param() + " characters")
			}
			if err.Tag() == "ascii" {
				return errors.New(field + " must be ascii characters")
			}
		}
		return err
	}
	return nil 
}