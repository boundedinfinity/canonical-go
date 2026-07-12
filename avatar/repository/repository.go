package repository

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type Repository[T any] interface {
	// Get retrieves an item by its [ider.Id].
	//
	// Returns the an item pointer and nil error if found.
	// Otherwise returns nil and an error.
	Get(ider.Id) (*T, error)
	// GetOk retrieves an item by its [ider.Id]. Returns the item and true if found, otherwise returns nil and false.
	GetOk(ider.Id) (*T, bool)
	// Find retrieves items that match the given term. Returns a slice of matching items and nil error.
	Find(term string) ([]*T, error)

	Match(ider.Id) ([]*T, error)
}
