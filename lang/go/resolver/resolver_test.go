package resolver_test

import (
	"errors"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

var (
	errNotFound = errors.New("not found")
	known       = id.MustParse("b12e7628-35b0-4ae5-b2d9-0f34f4828229")
	unknown     = id.MustParse("0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c")

	lookup = resolver.Func[id.ID, string](func(i id.ID) (string, error) {
		if i == known {
			return "found", nil
		}

		return "", errNotFound
	})
)

func Test_One(t *testing.T) {
	tcs := []struct {
		name     string
		resolver resolver.Resolver[id.ID, string]
		id       id.ID
		expected string
		errs     []error
	}{
		{name: "found", resolver: lookup, id: known, expected: "found"},
		{name: "not found", resolver: lookup, id: unknown, errs: []error{resolver.ErrResolve, errNotFound}},
		{name: "nil resolver", resolver: nil, id: known, errs: []error{resolver.ErrResolve, resolver.ErrNoResolver}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := resolver.One(tc.resolver, tc.id)

			for _, e := range tc.errs {
				assert.ErrorIs(tt, err, e)
			}

			if len(tc.errs) == 0 {
				assert.NoError(tt, err)
			}

			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_All(t *testing.T) {
	tcs := []struct {
		name     string
		ids      []id.ID
		expected []string
		errs     []error
	}{
		{name: "all found", ids: []id.ID{known, known}, expected: []string{"found", "found"}},
		{name: "empty", ids: []id.ID{}, expected: []string{}},
		{name: "one missing", ids: []id.ID{known, unknown}, errs: []error{resolver.ErrResolve, errNotFound}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := resolver.All(lookup, tc.ids)

			for _, e := range tc.errs {
				assert.ErrorIs(tt, err, e)
			}

			if len(tc.errs) == 0 {
				assert.NoError(tt, err)
			}

			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_All_reports_every_failure(t *testing.T) {
	_, err := resolver.All(lookup, []id.ID{unknown, known, unknown})

	assert.ErrorContains(t, err, "[0]")
	assert.ErrorContains(t, err, "[2]")
	assert.NotContains(t, err.Error(), "[1]")
}

func Test_Option(t *testing.T) {
	actual, err := resolver.Option(lookup, optioner.Some(known))
	assert.NoError(t, err)
	assert.Equal(t, optioner.Some("found"), actual)

	actual, err = resolver.Option[id.ID, string](nil, optioner.None[id.ID]())
	assert.NoError(t, err)
	assert.Equal(t, optioner.None[string](), actual)

	actual, err = resolver.Option(lookup, optioner.Some(unknown))
	assert.ErrorIs(t, err, resolver.ErrResolve)
	assert.Equal(t, optioner.None[string](), actual)
}

func Test_OptionAll(t *testing.T) {
	actual, err := resolver.OptionAll(lookup, optioner.Some([]id.ID{known}))
	assert.NoError(t, err)
	assert.Equal(t, optioner.Some([]string{"found"}), actual)

	actual, err = resolver.OptionAll[id.ID, string](nil, optioner.None[[]id.ID]())
	assert.NoError(t, err)
	assert.Equal(t, optioner.None[[]string](), actual)

	actual, err = resolver.OptionAll(lookup, optioner.Some([]id.ID{unknown}))
	assert.ErrorIs(t, err, resolver.ErrResolve)
	assert.Equal(t, optioner.None[[]string](), actual)
}
