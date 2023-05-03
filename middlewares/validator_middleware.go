package middlewares

import (
	"errors"
	"strings"

	"github.com/go-playground/validator"
)

type ValidatorMiddleware struct {
	Validator *validator.Validate
}

func NewValidatorMiddleware() *ValidatorMiddleware {
	return &ValidatorMiddleware{Validator: validator.New()}
}

func (v *ValidatorMiddleware) Validate(i interface{}) error {	
	if err := v.Validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			if err.Tag() == "required" {
				return errors.New(field + " is required")
			}
			if err.Tag() == "email" {
				return errors.New(field + " is not valid email")
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