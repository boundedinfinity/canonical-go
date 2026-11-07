package validation

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	FloatValidationError = errors.New("float validation")
)

func Float[T ~float32 | ~float64](name string) *floatValidator[T] {
	return &floatValidator[T]{
		name:    name,
		rootErr: IntegerValidationError,
	}
}

var _ Validator = &integerValidator[int]{}

type floatValidator[T ~float32 | ~float64] struct {
	validator[T]
}

func (this *floatValidator[T]) Add(fn func(T) error) *floatValidator[T] {
	return this.add(this, fn)
}

func (this *floatValidator[T]) Min(t T) *floatValidator[T] {
	validation := func(v T) error {
		if v < t {
			return fmt.Errorf("%f less than minimium %f", v, t)
		}

		return nil
	}

	return this.Add(validation)
}

func (this *floatValidator[T]) Max(t T) *floatValidator[T] {
	this.validations = append(this.validations, func(v T) error {
		if v > t {
			return fmt.Errorf("%f greater than maximum %f", v, t)
		}

		return nil
	})

	return this
}

func (this *floatValidator[T]) Positive() *floatValidator[T] {
	return this.Min(0)
}

func (this *floatValidator[T]) Negative() *floatValidator[T] {
	return this.Max(0)
}

func (this *floatValidator[T]) OneOf(ts ...T) *floatValidator[T] {
	validation := func(v T) error {
		if !slices.Contains(ts, v) {
			strs := make([]string, len(ts))

			for i := range ts {
				strs[i] = fmt.Sprintf("%f", ts[i])
			}

			return fmt.Errorf("%f is not one of %s", v, strings.Join(strs, ", "))
		}

		return nil
	}

	return this.Add(validation)
}

func (this *floatValidator[T]) NoneOf(ts ...T) *floatValidator[T] {
	validation := func(v T) error {
		if slices.Contains(ts, v) {
			strs := make([]string, len(ts))

			for i := range ts {
				strs[i] = fmt.Sprintf("%f", ts[i])
			}

			return fmt.Errorf("%f is one of %s", v, strings.Join(strs, ", "))
		}

		return nil
	}

	return this.Add(validation)
}
