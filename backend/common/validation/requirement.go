package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationErrors struct {
	Errors []string `json:"errors"`
}

// Alphebetical validation.
func alphabeticalValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	return regexp.MustCompile("^[a-zA-Z_ ]+$").MatchString(input)
}

// Alphanumeric validation.
func alphanumericalValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	return regexp.MustCompile("^[a-zA-Z0-9_ ]+$").MatchString(input)
}

// Alphanumspecial validation.
func alphanumspecialValidation(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	return regexp.MustCompile("^[a-zA-Z0-9 _\\-/:;()$&@\".,?!'\\[\\]{}#%^*+=|~<>·]+$").MatchString(input)
}

// Validates the request payload based on the struct's validate tags.
func ValidateRequest(requestPayload any) error {
	validate := validator.New()
	validate.RegisterValidation("alphabetical", alphabeticalValidation)
	validate.RegisterValidation("alphanumerical", alphanumericalValidation)
	validate.RegisterValidation("alphanumspecial", alphanumspecialValidation)

	err := validate.Struct(requestPayload)
	var param string
	if err != nil {
		var errs validationErrors
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid json")
		}
		for _, err := range err.(validator.ValidationErrors) {
			param = fmt.Sprintf("%s: %s", err.Tag(), err.Param())
			if err.Param() == "" {
				param = err.Tag()
			}
			errs.Errors = append(errs.Errors, fmt.Sprintf("field: %s, expected %s", err.Field(), param))
		}
		return errors.New(strings.Join(errs.Errors, ", "))
	}
	return nil
}
