package canonical

import "github.com/boundedinfinity/canonical-go/idiomatic/ider"

type OneToOne[T any] struct {
	loaded bool
	Id     ider.Id
	Value  T
}

type OneToMany[T any] struct {
	loaded bool
	Ids    []ider.Id
	Values []T
}

type ManyToMany[T any] struct {
	loaded bool
	Ids    []ider.Id
	Values []T
}
