package validation

import (
	"errors"
	"fmt"
)

type ValidationManager struct {
}

type Validator interface {
	Validate(a any) error
}

var _ Validator = &validator[any]{}

type validator[T any] struct {
	name        string
	validations []func(T) error
	rootErr     error
}

func (this *validator[T]) add[R any](r R, fn func(T) error) R {
	this.validations = append(this.validations, fn)
	return r
}

func (this *validator[T]) Name(v string) *validator[T] {
	this.name = v
	return this
}

func (this *validator[T]) Validate(a any) error {
	t, ok := a.(T)

	if !ok {
		var t T
		panic(fmt.Errorf("%T is not of type %T", a, t))
	}

	var errs []error

	for _, fn := range this.validations {
		if err := fn(t); err != nil {
			if this.name != "" {
				err = fmt.Errorf("%s %w", this.name, err)
			}
			errs = append(errs, fmt.Errorf("%w : %w", this.rootErr, err))
		}
	}

	return errors.Join(errs...)
}
