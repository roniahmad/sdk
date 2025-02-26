package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(data interface{}) error {
	validate := validator.New()

	errs := validate.Struct(data)
	if errs != nil {
		var errMsgs = make([]string, 0)
		for _, err := range errs.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s is '%s'", err.Field(), err.Tag()))
		}
		return errors.New(strings.Join(errMsgs, " and "))
	}

	return nil
}
