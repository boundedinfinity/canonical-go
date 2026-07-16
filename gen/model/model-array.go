package model

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	o "github.com/boundedinfinity/go-commoner/functional/optioner"
)

var _ Model = &ArrayModel{}

var (
	ErrArray             = errorer.New("array model error")
	ErrArrayValidation   = errorer.New("validation error")
	errArrayValidationFn = errorer.Func(ErrArray, ErrArrayValidation)
)

type ArrayModel struct {
	Name        string           `json:"name"`
	Description o.Option[string] `json:"description"`
	Required    o.Option[bool]   `json:"required"`
	Item        Model            `json:"item"`
	Min         o.Option[int]    `json:"min"`
	Max         o.Option[int]    `json:"max"`
}

func (_ ArrayModel) Kind() Kind {
	return Kinds.Array
}

func (this ArrayModel) Validate() error {
	errs := []error{}

	if this.Min.Defined() && this.Max.Defined() {
		if this.Min.Get() > this.Max.Get() {
			errs = append(errs, errArrayValidationFn(
				"min %v is greater than max %v",
				this.Min.Get(), this.Max.Get(),
			))
		}
	}

	return errorer.Join(errs...)
}

func Array(name string, item Model) *arrayBuilder {
	return &arrayBuilder{model: &ArrayModel{Name: name, Item: item}}
}

type arrayBuilder struct {
	model *ArrayModel
}

func (this *arrayBuilder) Kind() Kind {
	return this.model.Kind()
}

func (this arrayBuilder) Value() Model {
	return this.model
}

func (this *arrayBuilder) Name(name string) *arrayBuilder {
	this.model.Name = name
	return this
}

func (this *arrayBuilder) Description(description string) *arrayBuilder {
	this.model.Description = o.Some(description)
	return this
}

func (this *arrayBuilder) Item(item Model) *arrayBuilder {
	this.model.Item = item
	return this
}

func (this *arrayBuilder) Min(min int) *arrayBuilder {
	this.model.Min = o.Some(min)
	return this
}

func (this *arrayBuilder) Max(max int) *arrayBuilder {
	this.model.Max = o.Some(max)
	return this
}

func (this *arrayBuilder) Required(required bool) *arrayBuilder {
	this.model.Required = o.Some(required)
	return this
}
