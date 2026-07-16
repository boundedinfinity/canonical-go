package model

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	o "github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/boundedinfinity/go-commoner/idiomatic/mather"
)

var _ Model = &IntegerModel[int]{}

var (
	ErrInteger             = errorer.New("integer model error")
	ErrIntegerValidation   = errorer.New("validation error")
	errIntegerValidationFn = errorer.Func(ErrInteger, ErrIntegerValidation)
)

type IntegerModel[T mather.Integer] struct {
	Name        string           `json:"name"`
	Description o.Option[string] `json:"description"`
	Required    o.Option[bool]   `json:"required"`
	Min         o.Option[T]      `json:"min"`
	Max         o.Option[T]      `json:"max"`
	MultipleOf  o.Option[T]      `json:"multiple-of"`
	Positive    o.Option[bool]   `json:"positive"`
	Negative    o.Option[bool]   `json:"negative"`
	OneOf       o.Option[[]T]    `json:"one-of"`
}

// TODO: implement noneOf

func (_ IntegerModel[T]) Kind() Kind {
	return Kinds.Integer
}

func (this IntegerModel[T]) Validate() error {
	errs := []error{}

	if this.Positive.Defined() && this.Negative.Defined() {
		if this.Positive.Get() && this.Negative.Get() {
			errs = append(errs, errIntegerValidationFn("positive and negative are mutually exclusive"))
		}
	}

	if this.Min.Defined() && this.Max.Defined() {
		if this.Min.Get() > this.Max.Get() {
			errs = append(errs, errIntegerValidationFn(
				"min %v is greater than max %v",
				this.Min.Get(), this.Max.Get(),
			))
		}
	}

	if this.OneOf.Defined() {
		lessThanMin := []T{}
		greaterThanMax := []T{}

		for _, v := range this.OneOf.Get() {
			if this.Min.Defined() && v < this.Min.Get() {
				lessThanMin = append(lessThanMin, v)
			}

			if this.Max.Defined() && v > this.Max.Get() {
				greaterThanMax = append(greaterThanMax, v)
			}
		}

		if len(lessThanMin) > 0 {
			errs = append(errs, errIntegerValidationFn(
				"oneOf contains values less than min of %v: %v",
				this.Min.Get(), lessThanMin,
			))
		}

		if len(greaterThanMax) > 0 {
			errs = append(errs, errIntegerValidationFn(
				"oneOf contains values greater than max of %v: %v",
				this.Max.Get(), greaterThanMax,
			))
		}
	}

	return errorer.Join(errs...)
}

func Integer[T mather.Integer](name string) *integerBuilder[T] {
	return &integerBuilder[T]{model: &IntegerModel[T]{Name: name}}
}

type integerBuilder[T mather.Integer] struct {
	model *IntegerModel[T]
}

func (this *integerBuilder[T]) Kind() Kind {
	return this.model.Kind()
}

func (this integerBuilder[T]) Value() *IntegerModel[T] {
	return this.model
}

func (this *integerBuilder[T]) Name(name string) *integerBuilder[T] {
	this.model.Name = name
	return this
}

func (this *integerBuilder[T]) Description(description string) *integerBuilder[T] {
	this.model.Description = o.Some(description)
	return this
}

func (this *integerBuilder[T]) Min(min T) *integerBuilder[T] {
	this.model.Min = o.Some(min)
	return this
}

func (this *integerBuilder[T]) Max(max T) *integerBuilder[T] {
	this.model.Max = o.Some(max)
	return this
}

func (this *integerBuilder[T]) MultipleOf(multipleOf T) *integerBuilder[T] {
	this.model.MultipleOf = o.Some(multipleOf)
	return this
}

func (this *integerBuilder[T]) Positive(positive bool) *integerBuilder[T] {
	this.model.Positive = o.Some(positive)
	return this
}

func (this *integerBuilder[T]) Negative(negative bool) *integerBuilder[T] {
	this.model.Negative = o.Some(negative)
	return this
}

func (this *integerBuilder[T]) OneOf(oneOf []T) *integerBuilder[T] {
	if len(oneOf) > 0 {
		this.model.OneOf = o.Some(oneOf)
	}

	return this
}
