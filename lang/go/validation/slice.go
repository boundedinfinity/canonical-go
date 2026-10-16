package validation

import (
	"errors"
	"fmt"
)

var (
	SliceValidationError = errors.New("slice validation")
)

func Slice[T any](name string) *sliceValidator[T] {
	return &sliceValidator[T]{
		name: name,
	}
}

type sliceValidator[T any] struct {
	name        string
	validations []func([]T) error
}

func (this sliceValidator[T]) Validate(a any) error {
	v, ok := a.([]T)

	if !ok {
		panic(fmt.Errorf("%T is not of type %T", a, v))
	}

	var errs []error

	for _, fn := range this.validations {
		if err := fn(v); err != nil {
			if this.name != "" {
				err = fmt.Errorf("%s %w", this.name, err)
			}
			errs = append(errs, fmt.Errorf("%w : %w", SliceValidationError, err))
		}
	}

	return errors.Join(errs...)
}

func (this *sliceValidator[T]) Name(v string) *sliceValidator[T] {
	this.name = v
	return this
}

func (this *sliceValidator[T]) Add(fn func([]T) error) *sliceValidator[T] {
	this.validations = append(this.validations, fn)
	return this
}

func (this *sliceValidator[T]) Required() *sliceValidator[T] {
	validation := func(v []T) error {
		if len(v) == 0 {
			return fmt.Errorf("required")
		}

		return nil
	}

	return this.Add(validation)
}

func (this *sliceValidator[T]) Min(t int) *sliceValidator[T] {
	validation := func(v []T) error {
		l := len(v)

		if l < t {
			return fmt.Errorf("length %d less than minimium %d", l, t)
		}

		return nil
	}

	return this.Add(validation)
}

func (this *sliceValidator[T]) Max(t int) *sliceValidator[T] {
	validation := func(v []T) error {
		l := len(v)

		if l > t {
			return fmt.Errorf("length %d greater than maximum %d", l, t)
		}

		return nil
	}

	return this.Add(validation)
}

// Each runs the element validator on every element, reporting failures by index.
func (this *sliceValidator[T]) Each(validator Validator) *sliceValidator[T] {
	validation := func(v []T) error {
		var errs []error

		for i, e := range v {
			if err := validator.Validate(e); err != nil {
				errs = append(errs, fmt.Errorf("[%d] %w", i, err))
			}
		}

		return errors.Join(errs...)
	}

	return this.Add(validation)
}
