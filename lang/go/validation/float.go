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

type floatValidator[T ~float32 | ~float64] struct {
	name        string
	validations []func(T) error
}

func (this floatValidator[T]) Validate(a any) error {
	var v T
	var ok bool

	if v, ok = a.(T); !ok {
		panic(fmt.Errorf("%T is not of type %T", a, v))
	}

	var errs []error

	for _, fn := range this.validations {
		if err := fn(v); err != nil {
			if this.name != "" {
				err = fmt.Errorf("%s %w", this.name, err)
			}

			errs = append(errs, fmt.Errorf("%w : %w", FloatValidationError, err))
		}
	}

	return errors.Join(errs...)
}

func (this *floatValidator[T]) Name(v string) *floatValidator[T] {
	this.name = v
	return this
}

func (this *floatValidator[T]) Add(fn func(T) error) *floatValidator[T] {
	this.validations = append(this.validations, fn)
	return this
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
