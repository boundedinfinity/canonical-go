package validation

import (
	"errors"
	"fmt"
)

type Validator interface {
	Name() string
	Validate(a any) error
}

type ValidatationManager struct {
}

func runValidate[T any](name string, a any, validations []func(T) error) error {
	var v T
	var ok bool

	if v, ok = a.(T); !ok {
		panic(fmt.Errorf("%T is not of type %T", a, v))
	}

	var errs []error

	for _, fn := range validations {
		if err := fn(v); err != nil {
			if name != "" {
				err = fmt.Errorf("%s %w", name, err)
			}

			errs = append(errs, fmt.Errorf("%w : %w", IntegerValidationError, err))
		}
	}

	return errors.Join(errs...)
}
