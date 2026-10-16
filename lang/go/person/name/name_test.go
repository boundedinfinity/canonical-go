package name_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/canonical-go/lang/go/person/name"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	dr := affix.MustNewPrefix("Dr.", "Doctor", affix.MustNewPrefixCategory("Honorific", ""))
	jr := affix.MustNewSuffix("Jr.", "Junior", affix.MustNewSuffixCategory("Generational", ""))

	tcs := []struct {
		name   string
		prefix []affix.Prefix
		first  []string
		middle []string
		last   []string
		suffix []affix.Suffix
		errs   []error
	}{
		{name: "first only", first: []string{"Ada"}},
		{name: "middle only", middle: []string{"Augusta"}},
		{name: "last only", last: []string{"Lovelace"}},
		{name: "all parts", first: []string{"Ada"}, middle: []string{"Augusta"}, last: []string{"King", "Lovelace"}},
		{name: "with prefix", prefix: []affix.Prefix{dr}, last: []string{"Lovelace"}},
		{name: "multiple prefixes", prefix: []affix.Prefix{dr, dr}, last: []string{"Lovelace"}},
		{name: "with suffix", last: []string{"King"}, suffix: []affix.Suffix{jr}},
		{name: "multiple suffixes", last: []string{"King"}, suffix: []affix.Suffix{jr, jr}},
		{name: "part min", first: []string{"A"}},
		{name: "part max", first: []string{strings.Repeat("a", 30)}},
		{name: "part max counts runes", first: []string{strings.Repeat("é", 30)}},
		{name: "no parts", errs: []error{name.ErrRequired}},
		{name: "affixes only", prefix: []affix.Prefix{dr}, suffix: []affix.Suffix{jr}, errs: []error{name.ErrRequired}},
		{name: "empty prefix", prefix: []affix.Prefix{}, first: []string{"Ada"}},
		{name: "empty suffix", first: []string{"Ada"}, suffix: []affix.Suffix{}},
		{name: "zero prefix", prefix: []affix.Prefix{{}}, first: []string{"Ada"}, errs: []error{name.ErrPrefix, affix.ErrName, affix.ErrCategory, id.ErrInvalid}},
		{name: "invalid suffix", first: []string{"Ada"}, suffix: []affix.Suffix{jr, {Id: id.New(), Name: "Jr.", Description: optioner.Some("a"), Category: jr.Category}}, errs: []error{name.ErrSuffix, affix.ErrDescription}},
		{name: "empty first", first: []string{}, errs: []error{name.ErrRequired}},
		{name: "empty string", first: []string{""}, errs: []error{name.ErrFirst}},
		{name: "part too long", last: []string{strings.Repeat("a", 31)}, errs: []error{name.ErrLast}},
		{name: "bad element after good", first: []string{"Ada", ""}, errs: []error{name.ErrFirst}},
		{
			name:   "reports every failure",
			first:  []string{""},
			middle: []string{},
			last:   []string{strings.Repeat("a", 31)},
			errs:   []error{name.ErrFirst, name.ErrLast},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := name.New(tc.prefix, tc.first, tc.middle, tc.last, tc.suffix)

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
			assert.Equal(tt, optioner.OfSlice(tc.prefix), actual.Prefix)
			assert.Equal(tt, tc.first, actual.First.Get())
			assert.Equal(tt, tc.middle, actual.Middle.Get())
			assert.Equal(tt, tc.last, actual.Last.Get())
			assert.Equal(tt, optioner.OfSlice(tc.suffix), actual.Suffix)
		})
	}
}

func Test_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawPrefixId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"
	const rawSuffixId = "5c3a9e1d-7b2f-4e6a-8d0c-1f4b6a2e9c7d"

	const rawCategoryId = "9a7e4c2b-1d3f-4b5a-8e6c-2f0d9b8a7c6e"

	dr := affix.Prefix{
		Id:          id.MustParse(rawPrefixId),
		Name:        "Dr.",
		Description: optioner.Some("Doctor"),
		Category:    affix.PrefixCategory{Id: id.MustParse(rawCategoryId), Name: "Honorific"},
	}
	jr := affix.Suffix{
		Id:          id.MustParse(rawSuffixId),
		Name:        "Jr.",
		Description: optioner.Some("Junior"),
		Category:    affix.SuffixCategory{Id: id.MustParse(rawCategoryId), Name: "Generational"},
	}

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
				[]affix.Prefix{dr},
				[]string{"Ada"},
				[]string{"Augusta"},
				[]string{"King", "Lovelace"},
				[]affix.Suffix{jr},
			),
			expected: `{"id":"` + rawId + `",` +
				`"prefix":[{"id":"` + rawPrefixId + `","name":"Dr.","description":"Doctor","category":{"id":"` + rawCategoryId + `","name":"Honorific"}}],` +
				`"first":["Ada"],"middle":["Augusta"],"last":["King","Lovelace"],` +
				`"suffix":[{"id":"` + rawSuffixId + `","name":"Jr.","description":"Junior","category":{"id":"` + rawCategoryId + `","name":"Generational"}}]}`,
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
		{name: "empty prefix", input: `{"id":` + rawId + `,"prefix":[],"first":["Ada"]}`, err: name.ErrPrefix},
		{name: "empty suffix", input: `{"id":` + rawId + `,"first":["Ada"],"suffix":[]}`, err: name.ErrSuffix},
		{name: "invalid prefix", input: `{"id":` + rawId + `,"prefix":[{"id":` + rawId + `,"name":"","category":{"id":` + rawId + `,"name":"Honorific"}}],"first":["Ada"]}`, err: affix.ErrName},
		{name: "prefix without category", input: `{"id":` + rawId + `,"prefix":[{"id":` + rawId + `,"name":"Dr."}],"first":["Ada"]}`, err: affix.ErrCategory},
		{name: "invalid suffix", input: `{"id":` + rawId + `,"first":["Ada"],"suffix":[{"id":"nope"}]}`, err: id.ErrInvalid},
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
	expected := name.MustNew(nil, []string{"Ada"}, nil, nil, nil)
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
