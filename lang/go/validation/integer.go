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

func Integer[T ~int | ~int32 | ~int64](name string) *integerValidator[T] {
	return &integerValidator[T]{
		name:    name,
		rootErr: IntegerValidationError,
	}
}

var _ Validator = &integerValidator[int]{}

type integerValidator[T ~int | ~int32 | ~int64] struct {
	validator[T]
}

func (this *integerValidator[T]) Add(fn func(T) error) *integerValidator[T] {
	return this.add(this, fn)
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
