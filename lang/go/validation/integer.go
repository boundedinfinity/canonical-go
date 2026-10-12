package validation

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	IntegerValidationError = errors.New("string validation")
)

type integerValidator[T ~int | ~int32 | ~int64] struct {
	name        string
	validations []func(T) error
}

func (this integerValidator[T]) Validate(a any) error {
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
			errs = append(errs, fmt.Errorf("%w : %w", IntegerValidationError, err))
		}
	}

	return errors.Join(errs...)
}

func (this *integerValidator[T]) Name(v string) *integerValidator[T] {
	this.name = v
	return this
}

func (this *integerValidator[T]) Add(fn func(T) error) *integerValidator[T] {
	this.validations = append(this.validations, fn)
	return this
}

func (this *integerValidator[T]) Min(t T) *integerValidator[T] {
	validation := func(v T) error {
		if v < t {
			return fmt.Errorf("%d less than minimium %d", v, t)
		}

		return nil
	}

	return this.Add(validation)
}

func (this *integerValidator[T]) Max(t T) *integerValidator[T] {
	this.validations = append(this.validations, func(v T) error {
		if v > t {
			return fmt.Errorf("%d greater than maximum %d", v, t)
		}

		return nil
	})

	return this
}

func (this *integerValidator[T]) Positive() *integerValidator[T] {
	return this.Min(0)
}

func (this *integerValidator[T]) Negative() *integerValidator[T] {
	return this.Max(0)
}

func (this *integerValidator[T]) MultipleOf(t T) *integerValidator[T] {
	this.validations = append(this.validations, func(v T) error {
		if v%t != 0 {
			return fmt.Errorf("%d not mulitple of %d", v, t)
		}

		return nil
	})

	return this
}

func (this *integerValidator[T]) OneOf(ts ...T) *integerValidator[T] {
	validation := func(v T) error {
		if !slices.Contains(ts, v) {
			strs := make([]string, len(ts))

			for i := range ts {
				strs[i] = fmt.Sprintf("%d", ts[i])
			}

			return fmt.Errorf("%d is not one of %s", v, strings.Join(strs, ", "))
		}

		return nil
	}

	return this.Add(validation)
}

func (this *integerValidator[T]) NoneOf(ts ...T) *integerValidator[T] {
	validation := func(v T) error {
		if slices.Contains(ts, v) {
			strs := make([]string, len(ts))

			for i := range ts {
				strs[i] = fmt.Sprintf("%d", ts[i])
			}

			return fmt.Errorf("%d is one of %s", v, strings.Join(strs, ", "))
		}

		return nil
	}

	return this.Add(validation)
}
