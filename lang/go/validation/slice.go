package validation

import (
	"errors"
	"fmt"
)

var (
	SliceValidationError = errors.New("float validation")
)

type sliceValidator[T any] struct {
	name        string
	validations []func([]T) error
}

func (this *sliceValidator[T]) Add(fn func([]T) error) *sliceValidator[T] {
	this.validations = append(this.validations, fn)
	return this
}

func (this *sliceValidator[T]) Required() *sliceValidator[T] {
	validation := func(v []T) error {
		if v == nil {
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

		if l < t {
			return fmt.Errorf("length %d greater than maximum %d", l, t)
		}

		return nil
	}

	return this.Add(validation)
}
