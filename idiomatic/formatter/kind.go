package formater

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
)

var (
	ErrKind   = errorer.New("formatter kind")
	errKindfn = errorer.Func(ErrKind)
)

type Kind string

func (this Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(this)
}

func (this *Kind) UnmarshalJSON(data []byte) error {
	var term string

	if err := json.Unmarshal(data, this); err != nil {
		return errKindfn("%w", err)
	}

	if kind, ok := Kinds.ParseOk(term); ok {
		*this = kind
	} else {
		return errKindfn("%s is invalid", term)
	}

	return nil
}

var Kinds = kinds{
	Text: "canonical.formatter.text",
}

type kinds struct {
	Text Kind
}

func (this kinds) All() []Kind {
	return []Kind{
		this.Text,
	}
}

func (this kinds) ParseOk(term string) (Kind, bool) {
	kterm := Kind(term)
	return slicer.FindFn(func(kind Kind) bool {
		return kterm == kind
	})
}
