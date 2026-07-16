package model

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	o "github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/boundedinfinity/go-commoner/idiomatic/mather"
)

var _ Model = &FloatModel[float32]{}

var (
	ErrFloat             = errorer.New("float model error")
	ErrFloatValidation   = errorer.New("validation error")
	errFloatValidationFn = errorer.Func(ErrFloat, ErrFloatValidation)
)

type FloatModel[T mather.Float] struct {
	Name        string
	Description o.Option[string]
	Required    o.Option[bool]
	Epsilon     o.Option[float64]
	Min         o.Option[float64]
	Max         o.Option[float64]
	Positive    o.Option[bool]
	Negative    o.Option[bool]
}

// TODO: Implement float comparison with epsilon
// check IEEE 754

func (_ FloatModel[T]) Kind() Kind {
	return Kinds.Float
}

func (this FloatModel[T]) Validate() error {
	if this.Positive.Defined() && this.Negative.Defined() {
		if this.Positive.Get() && this.Negative.Get() {
			return errFloatValidationFn("positive and negative are mutually exclusive")
		}
	}

	if this.Min.Defined() && this.Max.Defined() {
		if this.Min.Get() > this.Max.Get() {
			return errFloatValidationFn("min is greater than max")
		}
	}

	return nil
}

func Float[T mather.Float](name string) *floatBuilder[T] {
	return &floatBuilder[T]{model: &FloatModel[T]{Name: name}}
}

type floatBuilder[T mather.Float] struct {
	model *FloatModel[T]
}

func (this *floatBuilder[T]) Kind() Kind {
	return this.model.Kind()
}

func (this floatBuilder[T]) Value() *FloatModel[T] {
	return this.model
}

func (this *floatBuilder[T]) Name(name string) *floatBuilder[T] {
	this.model.Name = name
	return this
}

func (this *floatBuilder[T]) Description(description string) *floatBuilder[T] {
	this.model.Description = o.Some(description)
	return this
}

func (this *floatBuilder[T]) Min(min float64) *floatBuilder[T] {
	this.model.Min = o.Some(min)
	return this
}

func (this *floatBuilder[T]) Max(max float64) *floatBuilder[T] {
	this.model.Max = o.Some(max)
	return this
}

func (this *floatBuilder[T]) Positive(positive bool) *floatBuilder[T] {
	this.model.Positive = o.Some(positive)
	return this
}

func (this *floatBuilder[T]) Negative(negative bool) *floatBuilder[T] {
	this.model.Negative = o.Some(negative)
	return this
}

func (this *floatBuilder[T]) Epsilon(epsilon float64) *floatBuilder[T] {
	this.model.Epsilon = o.Some(epsilon)
	return this
}

func (this *floatBuilder[T]) Required(required bool) *floatBuilder[T] {
	this.model.Required = o.Some(required)
	return this
}
