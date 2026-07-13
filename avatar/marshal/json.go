package marshal

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	ErrKindMarshaler      = errorer.New("kind marshaller error")
	errKindMarshalerFn    = errorer.Func(ErrKindMarshaler)
	ErrKindMarshal        = errorer.New("marshal error")
	errKindMarshallerFn   = errorer.Func(ErrKindMarshaler, ErrKindMarshal)
	ErrKindUnmarshal      = errorer.New("unmarshal error")
	errKindUnmarshallerFn = errorer.Func(ErrKindMarshaler, ErrKindUnmarshal)
)

type descriminator struct {
	Kind string `json:"kind"`
}

func Marshal[T any](kind string, v T) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		err = errKindMarshallerFn("%w", err)
		return nil, err
	}

	var vmap map[string]any
	if err := json.Unmarshal(data, &vmap); err != nil {
		err = errKindMarshallerFn("%w", err)
		return nil, err
	}

	vmap["kind"] = kind

	data, err = json.Marshal(vmap)
	if err != nil {
		err = errKindMarshallerFn("%w", err)
		return nil, err
	}

	return data, nil
}
