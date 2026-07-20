package model

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	o "github.com/boundedinfinity/go-commoner/functional/optioner"
)

var _ Model = &ObjectModel{}

var (
	ErrObject             = errorer.New("object model error")
	ErrObjectValidation   = errorer.New("validation error")
	errObjectValidationFn = errorer.Func(ErrObject, ErrObjectValidation)
)

type ObjectModel struct {
	Name        string           `json:"name"`
	Description o.Option[string] `json:"description"`
	Required    o.Option[bool]   `json:"required"`
	Properties  []Property       `json:"properties"`
}

func (_ ObjectModel) Kind() Kind {
	return Kinds.Object
}

type Property struct {
	Kind     Model
	Optional bool `json:"optional"`
}

func (this ObjectModel) Validate() error {
	errs := []error{}

	for _, property := range this.Properties {
		if err := property.Kind.Validate(); err != nil {
			var name string
			errs = append(errs, errObjectValidationFn("property %v : %w", name, err))
		}
	}

	return errorer.Join(errs...)
}

func Object(name string) *objectBuilder {
	return &objectBuilder{model: &ObjectModel{Name: name}}
}

type objectBuilder struct {
	model *ObjectModel
}
