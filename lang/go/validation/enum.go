package validation

import "errors"

var (
	EnumValidationError = errors.New("enum validation")
)

func Enum[T ~string](name string, items ...T) *enumValidator[T] {
	return &enumValidator[T]{
		name:    name,
		rootErr: StringValidationError,
	}
}

var _ Validator = &enumValidator[string]{}

type enumValidator[T ~string] struct {
	stringValidator[T]
}
