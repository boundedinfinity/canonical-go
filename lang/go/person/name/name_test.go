package name_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/name"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	tcs := []struct {
		name   string
		first  []string
		middle []string
		last   []string
		errs   []error
	}{
		{name: "first only", first: []string{"Ada"}},
		{name: "middle only", middle: []string{"Augusta"}},
		{name: "last only", last: []string{"Lovelace"}},
		{name: "all parts", first: []string{"Ada"}, middle: []string{"Augusta"}, last: []string{"King", "Lovelace"}},
		{name: "part min", first: []string{"A"}},
		{name: "part max", first: []string{strings.Repeat("a", 30)}},
		{name: "part max counts runes", first: []string{strings.Repeat("é", 30)}},
		{name: "no parts", errs: []error{name.ErrRequired}},
		{name: "empty first", first: []string{}, errs: []error{name.ErrFirst}},
		{name: "empty middle", middle: []string{}, last: []string{"Lovelace"}, errs: []error{name.ErrMiddle}},
		{name: "empty last", last: []string{}, first: []string{"Ada"}, errs: []error{name.ErrLast}},
		{name: "empty string", first: []string{""}, errs: []error{name.ErrFirst}},
		{name: "part too long", last: []string{strings.Repeat("a", 31)}, errs: []error{name.ErrLast}},
		{name: "bad element after good", first: []string{"Ada", ""}, errs: []error{name.ErrFirst}},
		{
			name:   "reports every failure",
			first:  []string{""},
			middle: []string{},
			last:   []string{strings.Repeat("a", 31)},
			errs:   []error{name.ErrFirst, name.ErrMiddle, name.ErrLast},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := name.New(tc.first, tc.middle, tc.last)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, name.ErrInvalid)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, name.Name{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.False(tt, actual.Id.IsZero())
			assert.Equal(tt, tc.first, actual.First.Get())
			assert.Equal(tt, tc.middle, actual.Middle.Get())
			assert.Equal(tt, tc.last, actual.Last.Get())
		})
	}
}

func Test_Validate_zero_id(t *testing.T) {
	err := name.MustNew([]string{"Ada"}, nil, nil).Validate()

	assert.ErrorIs(t, err, name.ErrInvalid)
	assert.ErrorIs(t, err, name.ErrId)
}

func Test_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	tcs := []struct {
		name     string
		input    name.Name
		expected string
	}{
		{
			name:     "absent parts are omitted",
			input:    name.Name{Id: id.MustParse(rawId), First: optioner.OfSlice([]string{"Ada"})},
			expected: `{"id":"` + rawId + `","first":["Ada"]}`,
		},
		{
			name: "all parts",
			input: name.MustCreate(
				id.MustParse(rawId),
				[]string{"Ada"},
				[]string{"Augusta"},
				[]string{"King", "Lovelace"},
			),
			expected: `{"id":"` + rawId + `","first":["Ada"],"middle":["Augusta"],"last":["King","Lovelace"]}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual name.Name
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"first":["Ada"]}`, err: name.ErrId},
		{name: "bad id", input: `{"id":"nope","first":["Ada"]}`, err: id.ErrInvalid},
		{name: "no parts", input: `{"id":` + rawId + `}`, err: name.ErrRequired},
		{name: "null parts", input: `{"id":` + rawId + `,"first":null}`, err: name.ErrRequired},
		{name: "empty array", input: `{"id":` + rawId + `,"first":[]}`, err: name.ErrFirst},
		{name: "empty string", input: `{"id":` + rawId + `,"last":[""]}`, err: name.ErrLast},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual name.Name
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, name.Name{}, actual)
		})
	}
}

func Test_UnmarshalJSON_null(t *testing.T) {
	expected := name.MustNew([]string{"Ada"}, nil, nil)
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
