// Package resolver looks up Top Level Models by key (an id.ID, or a natural
// key such as an ISO country code) when a DTO is converted back to its model.
//
// Source of truth: domain-model/global-notes/dto.md
package resolver

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	// ErrResolve is wrapped by every failure to resolve a key.
	ErrResolve = errorer.New("resolve")
	// ErrNoResolver is returned when a key must be resolved but the
	// Resolver is nil.
	ErrNoResolver = errorer.New("no resolver")

	errResolveFn = errorer.Func(ErrResolve)
)

// Resolver looks up a T by its key K.
type Resolver[K any, T any] interface {
	Resolve(K) (T, error)
}

// Func adapts an ordinary function to a Resolver.
type Func[K any, T any] func(K) (T, error)

// Resolve calls this(key).
func (this Func[K, T]) Resolve(key K) (T, error) {
	return this(key)
}

// One resolves a single key. A failure wraps ErrResolve and the resolver's
// error.
func One[K any, T any](resolver Resolver[K, T], key K) (T, error) {
	var zero T

	if resolver == nil {
		return zero, errResolveFn("%v: %w", key, ErrNoResolver)
	}

	v, err := resolver.Resolve(key)

	if err != nil {
		return zero, errResolveFn("%v: %w", key, err)
	}

	return v, nil
}

// All resolves every key in order, reporting all failures by index.
func All[K any, T any](resolver Resolver[K, T], keys []K) ([]T, error) {
	vs := make([]T, len(keys))
	var errs []error

	for i, key := range keys {
		v, err := One(resolver, key)

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

// Option resolves an optional key. None resolves to None without calling the
// resolver, so the resolver may be nil when the key is absent.
func Option[K any, T any](resolver Resolver[K, T], key optioner.Option[K]) (optioner.Option[T], error) {
	if key.Empty() {
		return optioner.None[T](), nil
	}

	v, err := One(resolver, key.Get())

	if err != nil {
		return optioner.None[T](), err
	}

	return optioner.Some(v), nil
}

// OptionAll resolves an optional list of keys. None resolves to None without
// calling the resolver, so the resolver may be nil when the list is absent.
func OptionAll[K any, T any](resolver Resolver[K, T], keys optioner.Option[[]K]) (optioner.Option[[]T], error) {
	if keys.Empty() {
		return optioner.None[[]T](), nil
	}

	vs, err := All(resolver, keys.Get())

	if err != nil {
		return optioner.None[[]T](), err
	}

	return optioner.Some(vs), nil
}
