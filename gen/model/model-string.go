package model

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	o "github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	ErrString             = errorer.New("string model error")
	ErrStringValidation   = errorer.New("validation error")
	errStringValidationFn = errorer.Func(ErrString, ErrStringValidation)
)

var _ Model = &StringModel{}

type StringModel struct {
	Name        string           `json:"name"`
	Description o.Option[string] `json:"description"`
	Required    o.Option[bool]   `json:"required"`
	Min         o.Option[int]    `json:"min"`
	Max         o.Option[int]    `json:"max"`
	CapitalCase o.Option[bool]   `json:"capital-case"`
	Uppercase   o.Option[bool]   `json:"uppercase"`
	Lowercase   o.Option[bool]   `json:"lowercase"`
	CamelCase   o.Option[bool]   `json:"camel-case"`
	SnakeCase   o.Option[bool]   `json:"snake-case"`
	PascalCase  o.Option[bool]   `json:"pascal-case"`
	KabaCase    o.Option[bool]   `json:"kaba-case"`
}

func (_ StringModel) Kind() Kind {
	return Kinds.String
}

func (this StringModel) Validate() error {
	errs := []error{}

	if this.Min.Defined() && this.Max.Defined() {
		if this.Min.Get() > this.Max.Get() {
			errs = append(errs, errStringValidationFn(
				"min %v is greater than max %v",
				this.Min.Get(), this.Max.Get(),
			))
		}

		if this.Min.Get() < 0 {
			errs = append(errs, errStringValidationFn(
				"min %v is less than 0",
				this.Min.Get(),
			))
		}

		if this.Max.Get() < 0 {
			errs = append(errs, errStringValidationFn(
				"max %v is less than 0",
				this.Max.Get(),
			))
		}
	}

	return errorer.Join(errs...)
}

func String[T ~string](name string) *stringBuilder[T] {
	return &stringBuilder[T]{model: &StringModel{Name: name}}
}

type stringBuilder[T ~string] struct {
	model *StringModel
}

func (this *stringBuilder[T]) Kind() Kind {
	return this.model.Kind()
}

func (this *stringBuilder[T]) Value() Model {
	return this.model
}

func (this *stringBuilder[T]) Name(name string) *stringBuilder[T] {
	this.model.Name = name
	return this
}

func (this *stringBuilder[T]) Description(description string) *stringBuilder[T] {
	this.model.Description = o.Some(description)
	return this
}

func (this *stringBuilder[T]) Min(min int) *stringBuilder[T] {
	this.model.Min = o.Some(min)
	return this
}

func (this *stringBuilder[T]) Max(max int) *stringBuilder[T] {
	this.model.Max = o.Some(max)
	return this
}

func (this *stringBuilder[T]) CapitalCase(capitalCase bool) *stringBuilder[T] {
	this.model.CapitalCase = o.Some(capitalCase)
	return this
}

func (this *stringBuilder[T]) Uppercase(uppercase bool) *stringBuilder[T] {
	this.model.Uppercase = o.Some(uppercase)
	return this
}

func (this *stringBuilder[T]) Lowercase(lowercase bool) *stringBuilder[T] {
	this.model.Lowercase = o.Some(lowercase)
	return this
}

func (this *stringBuilder[T]) CamelCase(camelCase bool) *stringBuilder[T] {
	this.model.CamelCase = o.Some(camelCase)
	return this
}

func (this *stringBuilder[T]) SnakeCase(snakeCase bool) *stringBuilder[T] {
	this.model.SnakeCase = o.Some(snakeCase)
	return this
}

func (this *stringBuilder[T]) PascalCase(pascalCase bool) *stringBuilder[T] {
	this.model.PascalCase = o.Some(pascalCase)
	return this
}

func (this *stringBuilder[T]) KabaCase(kabaCase bool) *stringBuilder[T] {
	this.model.KabaCase = o.Some(kabaCase)
	return this
}
