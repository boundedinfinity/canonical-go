package modeller

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/idiomatic/errorer"
)

type Kinder interface {
	Kind() string
}

var (
	ErrKindMarshaler      = errorer.New("kind marshaller error")
	errKindMarshalerFn    = errorer.Func(ErrKindMarshaler)
	ErrKindMarshal        = errorer.New("marshal error")
	errKindMarshallerFn   = errorer.Func(ErrKindMarshaler, ErrKindMarshal)
	ErrKindUnmarshal      = errorer.New("unmarshal error")
	errKindUnmarshallerFn = errorer.Func(ErrKindMarshaler, ErrKindUnmarshal)
)

func CanonicalJsonMarshal(kind Kinder) ([]byte, error) {
	data, err := json.Marshal(struct {
		Kind string
		Kinder
	}{
		Kind:   kind.Kind(),
		Kinder: kind,
	})

	if err != nil {
		err = errKindMarshallerFn("%w", err)
	}

	return data, err
}

type KindMarshaller struct {
	registry map[string]func(bytes []byte) error
}

type canonicalDescriminator struct {
	Kind string `json:"kind"`
}

func (this *KindMarshaller) Register(kind Kinder, fn func(bytes []byte) error) error {
	if _, ok := this.registry[kind.Kind()]; ok {
		return errKindMarshalerFn("kind '%s' already registered", kind)
	}

	this.registry[kind.Kind()] = fn

	return nil
}

func (this KindMarshaller) Unmarshal(data []byte) error {
	var descrim canonicalDescriminator

	if err := json.Unmarshal(data, &descrim); err != nil {
		return errKindUnmarshallerFn("%w", err)
	}

	if descrim.Kind == "" {
		return errKindUnmarshallerFn("kind not found in data")
	}

	fn, ok := this.registry[descrim.Kind]

	if !ok {
		return errKindUnmarshallerFn("kind '%s' not registered", descrim.Kind)
	}

	if err := fn(data); err != nil {
		return errKindUnmarshallerFn("%w", err)
	}

	return nil
}
