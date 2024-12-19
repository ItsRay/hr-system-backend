package errors

import (
	"errors"
	"fmt"
	"strings"
)

var ErrResourceNotFound = errors.New("resource not found")
var ErrInvalidInput = errors.New("invalid input")
var ErrStatusConflict = errors.New("status conflict")

func Combine(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	}
	if len(nonNilErrs) == 1 {
		return nonNilErrs[0]
	}

	var errStrings []string
	for _, err := range nonNilErrs {
		errStrings = append(errStrings, err.Error())
	}
	return fmt.Errorf("multiple errors: %s", strings.Join(errStrings, ", "))
}
