package repository

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/go-commoner/idiomatic/errorer"
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
)

func Memory[T any](
	name string,
	items []*T,
	byId func(id ider.Id) func(item *T) bool,
	byTerm func(term string) func(item *T) bool,
	byMatch func(id ider.Id) func(item *T) bool,
) Repository[T] {
	err := errorer.New(name + " repository error")

	return &MemoryManager[T]{
		Items:                items,
		ErrNotFoundFunc:      errorer.Func(err),
		ErrFoundMultipleFunc: errorer.Func(err),
		ById:                 byId,
		ByTerm:               byTerm,
		ByMatch:              byMatch,
	}
}

type MemoryManager[T any] struct {
	Items                []*T
	ErrNotFoundFunc      func(format string, a ...any) error
	ErrFoundMultipleFunc func(format string, a ...any) error
	ById                 func(id ider.Id) func(item *T) bool
	ByTerm               func(term string) func(item *T) bool
	ByMatch              func(id ider.Id) func(item *T) bool
}

func (this MemoryManager[T]) Get(id ider.Id) (*T, error) {
	results := slicer.FilterFunc(this.ById(id), this.Items...)

	switch len(results) {
	case 0:
		return nil, this.ErrNotFoundFunc("no item found for id: %s", id)
	case 1:
		return results[0], nil
	default:
		return nil, this.ErrFoundMultipleFunc("multiple items found for id: %s", id)
	}
}

func (this MemoryManager[T]) GetOk(id ider.Id) (*T, bool) {

	results := slicer.FilterFunc(this.ById(id), this.Items...)

	if len(results) == 0 {
		return nil, false
	}

	return results[0], true
}

func (this MemoryManager[T]) Find(term string) ([]*T, error) {
	return slicer.FilterFunc(this.ByTerm(term), this.Items...), nil
}

func (this MemoryManager[T]) Match(id ider.Id) ([]*T, error) {
	return slicer.FilterFunc(this.ByMatch(id), this.Items...), nil
}
