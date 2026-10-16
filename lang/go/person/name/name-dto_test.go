package name_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/canonical-go/lang/go/person/name"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

// lookup returns a Resolver over the given models, failing with errNotFound
// for any other ID.
func lookup[T any](errNotFound error, getId func(T) id.ID, models ...T) resolver.Resolver[T] {
	return resolver.Func[T](func(i id.ID) (T, error) {
		for _, m := range models {
			if getId(m) == i {
				return m, nil
			}
		}

		var zero T
		return zero, errNotFound
	})
}

func Test_NameDto_ToDto(t *testing.T) {
	dr := affix.MustNewPrefix("Dr.", "", affix.MustNewPrefixCategory("Honorific", ""))
	sir := affix.MustNewPrefix("Sir", "", affix.MustNewPrefixCategory("Honorific", ""))
	jr := affix.MustNewSuffix("Jr.", "", affix.MustNewSuffixCategory("Generational", ""))

	model := name.MustNew([]affix.Prefix{dr, sir}, []string{"Ada"}, nil, []string{"King"}, []affix.Suffix{jr})
	dto := model.ToDto()

	assert.Equal(t, model.Id, dto.Id)
	assert.Equal(t, optioner.Some([]id.ID{dr.Id, sir.Id}), dto.Prefix)
	assert.Equal(t, model.First, dto.First)
	assert.Equal(t, model.Middle, dto.Middle)
	assert.Equal(t, model.Last, dto.Last)
	assert.Equal(t, optioner.Some([]id.ID{jr.Id}), dto.Suffix)

	bare := name.MustNew(nil, []string{"Ada"}, nil, nil, nil).ToDto()
	assert.True(t, bare.Prefix.Empty())
	assert.True(t, bare.Suffix.Empty())
}

func Test_NameDto_ToModel(t *testing.T) {
	errNotFound := errors.New("not found")
	prefixId := func(p affix.Prefix) id.ID { return p.Id }
	suffixId := func(s affix.Suffix) id.ID { return s.Id }

	dr := affix.MustNewPrefix("Dr.", "", affix.MustNewPrefixCategory("Honorific", ""))
	jr := affix.MustNewSuffix("Jr.", "", affix.MustNewSuffixCategory("Generational", ""))
	badSuffix := affix.Suffix{Id: id.New(), Name: "Jr."}

	prefixes := lookup(errNotFound, prefixId, dr)
	suffixes := lookup(errNotFound, suffixId, jr, badSuffix)

	full := name.MustNew([]affix.Prefix{dr}, []string{"Ada"}, []string{"Augusta"}, []string{"King"}, []affix.Suffix{jr})
	bare := name.MustNew(nil, []string{"Ada"}, nil, nil, nil)

	tcs := []struct {
		name     string
		input    name.NameDto
		prefixes resolver.Resolver[affix.Prefix]
		suffixes resolver.Resolver[affix.Suffix]
		expected name.Name
		errs     []error
	}{
		{name: "round trip", input: full.ToDto(), prefixes: prefixes, suffixes: suffixes, expected: full},
		{name: "no references needs no resolvers", input: bare.ToDto(), expected: bare},
		{
			name:  "invalid dto",
			input: name.NameDto{Id: id.New(), Prefix: optioner.Some([]id.ID{})},
			errs:  []error{name.ErrInvalid, name.ErrRequired, name.ErrPrefix},
		},
		{
			name:     "unknown prefix and suffix are both reported",
			input:    name.NameDto{Id: id.New(), Prefix: optioner.Some([]id.ID{id.New()}), First: optioner.Some([]string{"Ada"}), Suffix: optioner.Some([]id.ID{jr.Id, id.New()})},
			prefixes: prefixes,
			suffixes: suffixes,
			errs:     []error{resolver.ErrResolve, errNotFound},
		},
		{
			name:     "missing resolver",
			input:    full.ToDto(),
			prefixes: prefixes,
			errs:     []error{resolver.ErrResolve, resolver.ErrNoResolver},
		},
		{
			name:     "resolved suffix is invalid",
			input:    name.NameDto{Id: id.New(), First: optioner.Some([]string{"Ada"}), Suffix: optioner.Some([]id.ID{badSuffix.Id})},
			suffixes: suffixes,
			errs:     []error{name.ErrInvalid, name.ErrSuffix, affix.ErrCategory},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := tc.input.ToModel(tc.prefixes, tc.suffixes)

			if len(tc.errs) > 0 {
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, name.Name{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_NameDto_ToModel_reports_prefix_and_suffix(t *testing.T) {
	errNotFound := errors.New("not found")
	none := lookup(errNotFound, func(p affix.Prefix) id.ID { return p.Id })
	noneSuffix := lookup(errNotFound, func(s affix.Suffix) id.ID { return s.Id })

	dto := name.NameDto{
		Id:     id.New(),
		Prefix: optioner.Some([]id.ID{id.New()}),
		First:  optioner.Some([]string{"Ada"}),
		Suffix: optioner.Some([]id.ID{id.New()}),
	}

	_, err := dto.ToModel(none, noneSuffix)

	assert.ErrorContains(t, err, "prefix")
	assert.ErrorContains(t, err, "suffix")
}

func Test_NameDto_Validate(t *testing.T) {
	validId := id.MustParse("b12e7628-35b0-4ae5-b2d9-0f34f4828229")
	ada := optioner.Some([]string{"Ada"})

	tcs := []struct {
		name  string
		input name.NameDto
		errs  []error
	}{
		{name: "first only", input: name.NameDto{Id: validId, First: ada}},
		{name: "with references", input: name.NameDto{Id: validId, Prefix: optioner.Some([]id.ID{validId}), First: ada, Suffix: optioner.Some([]id.ID{validId})}},
		{name: "zero id", input: name.NameDto{First: ada}, errs: []error{name.ErrId}},
		{name: "no parts", input: name.NameDto{Id: validId}, errs: []error{name.ErrRequired}},
		{name: "references only", input: name.NameDto{Id: validId, Prefix: optioner.Some([]id.ID{validId})}, errs: []error{name.ErrRequired}},
		{name: "empty first", input: name.NameDto{Id: validId, First: optioner.Some([]string{""})}, errs: []error{name.ErrFirst}},
		{name: "empty prefix", input: name.NameDto{Id: validId, Prefix: optioner.Some([]id.ID{}), First: ada}, errs: []error{name.ErrPrefix}},
		{name: "zero prefix id", input: name.NameDto{Id: validId, Prefix: optioner.Some([]id.ID{{}}), First: ada}, errs: []error{name.ErrPrefix, id.ErrInvalid}},
		{name: "empty suffix", input: name.NameDto{Id: validId, First: ada, Suffix: optioner.Some([]id.ID{})}, errs: []error{name.ErrSuffix}},
		{name: "zero suffix id", input: name.NameDto{Id: validId, First: ada, Suffix: optioner.Some([]id.ID{validId, {}})}, errs: []error{name.ErrSuffix, id.ErrInvalid}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, name.ErrInvalid)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				return
			}

			assert.NoError(tt, err)
		})
	}
}

func Test_NameDto_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawPrefixId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"
	const rawSuffixId = "5c3a9e1d-7b2f-4e6a-8d0c-1f4b6a2e9c7d"

	tcs := []struct {
		name     string
		input    name.NameDto
		expected string
	}{
		{
			name:     "absent parts are omitted",
			input:    name.NameDto{Id: id.MustParse(rawId), First: optioner.Some([]string{"Ada"})},
			expected: `{"id":"` + rawId + `","first":["Ada"]}`,
		},
		{
			name: "all parts",
			input: name.NameDto{
				Id:     id.MustParse(rawId),
				Prefix: optioner.Some([]id.ID{id.MustParse(rawPrefixId)}),
				First:  optioner.Some([]string{"Ada"}),
				Middle: optioner.Some([]string{"Augusta"}),
				Last:   optioner.Some([]string{"King", "Lovelace"}),
				Suffix: optioner.Some([]id.ID{id.MustParse(rawSuffixId)}),
			},
			expected: `{"id":"` + rawId + `","prefix":["` + rawPrefixId + `"],` +
				`"first":["Ada"],"middle":["Augusta"],"last":["King","Lovelace"],` +
				`"suffix":["` + rawSuffixId + `"]}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual name.NameDto
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_NameDto_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"first":["Ada"]}`, err: name.ErrId},
		{name: "no parts", input: `{"id":` + rawId + `}`, err: name.ErrRequired},
		{name: "empty prefix", input: `{"id":` + rawId + `,"prefix":[],"first":["Ada"]}`, err: name.ErrPrefix},
		{name: "bad prefix id", input: `{"id":` + rawId + `,"prefix":["nope"],"first":["Ada"]}`, err: id.ErrInvalid},
		{name: "embedded prefix object", input: `{"id":` + rawId + `,"prefix":[{"id":` + rawId + `}],"first":["Ada"]}`, err: id.ErrInvalid},
		{name: "empty suffix", input: `{"id":` + rawId + `,"first":["Ada"],"suffix":[]}`, err: name.ErrSuffix},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual name.NameDto
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, name.NameDto{}, actual)
		})
	}
}

func Test_NameDto_UnmarshalJSON_null(t *testing.T) {
	expected := name.MustNew(nil, []string{"Ada"}, nil, nil, nil).ToDto()
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
