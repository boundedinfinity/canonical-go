package model

import o "github.com/boundedinfinity/go-commoner/functional/optioner"

var _ Model = &BooleanModel{}

type BooleanModel struct {
	Name        string           `json:"name"`
	Description o.Option[string] `json:"description"`
}

func (_ BooleanModel) GetKind() Kind {
	return Kinds.Boolean
}

func (_ BooleanModel) Kind() Kind {
	return Kinds.Boolean
}

func (this BooleanModel) Validate() error {
	return nil
}

func Boolean(name string) *booleanBuilder {
	return &booleanBuilder{model: &BooleanModel{Name: name}}
}

type booleanBuilder struct {
	model *BooleanModel
}

func (this *booleanBuilder) Kind() Kind {
	return this.model.Kind()
}

func (this booleanBuilder) Value() Model {
	return this.model
}

func (this *booleanBuilder) Name(name string) *booleanBuilder {
	this.model.Name = name
	return this
}

func (this *booleanBuilder) Description(description string) *booleanBuilder {
	this.model.Description = o.Some(description)
	return this
}
