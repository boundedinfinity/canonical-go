package model

import o "github.com/boundedinfinity/go-commoner/functional/optioner"

type Model interface {
	Kind() Kind
	Validate() error
}

type models struct {
}

func (m *models) GetName(model Model) o.Option[string] {
	var name string

	switch m := model.(type) {
	case *ObjectModel:
		name = m.Name
	case *StringModel:
		name = m.Name
		// case *IntegerModel[int], *IntegerModel[int8], *IntegerModel[int16], *IntegerModel[int32], *IntegerModel[int64]:
		// 	name = m.Name

	}

	return o.OfZero(name)
}
