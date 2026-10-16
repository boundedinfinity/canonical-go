// Package resolver looks up Top Level Models by ID when a DTO is converted
// back to its model.
//
// Source of truth: domain-model/global-notes/dto.md
package resolver

import (
	"fmt"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	// ErrResolve is wrapped by every failure to resolve an ID.
	ErrResolve = errorer.New("resolve")
	// ErrNoResolver is returned when an ID must be resolved but the
	// Resolver is nil.
	ErrNoResolver = errorer.New("no resolver")

	errResolveFn = errorer.Func(ErrResolve)
)

// Resolver looks up a T by its ID.
type Resolver[T any] interface {
	Resolve(id.ID) (T, error)
}

// Func adapts an ordinary function to a Resolver.
type Func[T any] func(id.ID) (T, error)

// Resolve calls this(id).
func (this Func[T]) Resolve(id id.ID) (T, error) {
	return this(id)
}

// One resolves a single ID. A failure wraps ErrResolve and the resolver's error.
func One[T any](resolver Resolver[T], id id.ID) (T, error) {
	var zero T

	if resolver == nil {
		return zero, errResolveFn("%s: %w", id, ErrNoResolver)
	}

	v, err := resolver.Resolve(id)

	if err != nil {
		return zero, errResolveFn("%s: %w", id, err)
	}

	return v, nil
}

// All resolves every ID in order, reporting all failures by index.
func All[T any](resolver Resolver[T], ids []id.ID) ([]T, error) {
	vs := make([]T, len(ids))
	var errs []error

	for i, id := range ids {
		v, err := One(resolver, id)

		if err != nil {
			errs = append(errs, fmt.Errorf("[%d] %w", i, err))
			continue
		}

		vs[i] = v
	}

	if err := errorer.Join(errs...); err != nil {
		return nil, err
	}

	return vs, nil
}

// Option resolves an optional ID. None resolves to None without calling the
// resolver, so the resolver may be nil when the ID is absent.
func Option[T any](resolver Resolver[T], id optioner.Option[id.ID]) (optioner.Option[T], error) {
	if id.Empty() {
		return optioner.None[T](), nil
	}

	v, err := One(resolver, id.Get())

	if err != nil {
		return optioner.None[T](), err
	}

	return optioner.Some(v), nil
}

// OptionAll resolves an optional list of IDs. None resolves to None without
// calling the resolver, so the resolver may be nil when the list is absent.
func OptionAll[T any](resolver Resolver[T], ids optioner.Option[[]id.ID]) (optioner.Option[[]T], error) {
	if ids.Empty() {
		return optioner.None[[]T](), nil
	}

	vs, err := All(resolver, ids.Get())

	if err != nil {
		return optioner.None[[]T](), err
	}

	return optioner.Some(vs), nil
}
