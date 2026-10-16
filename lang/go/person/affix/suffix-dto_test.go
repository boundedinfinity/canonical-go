package affix_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_SuffixDto_ToDto(t *testing.T) {
	category := affix.MustNewSuffixCategory("Generational", "")
	model := affix.MustNewSuffix("Jr.", "Junior", category)

	dto := model.ToDto()

	assert.Equal(t, model.Affix, dto.Affix)
	assert.Equal(t, category.Id, dto.Category)
}

func Test_SuffixDto_ToModel(t *testing.T) {
	category := affix.MustNewSuffixCategory("Generational", "")
	invalidCategory := affix.SuffixCategory{Id: id.New(), Name: ""}
	errNotFound := errors.New("not found")

	categories := resolver.Func[affix.SuffixCategory](func(i id.ID) (affix.SuffixCategory, error) {
		switch i {
		case category.Id:
			return category, nil
		case invalidCategory.Id:
			return invalidCategory, nil
		}

		return affix.SuffixCategory{}, errNotFound
	})

	valid := affix.MustNewSuffix("Jr.", "Junior", category)

	tcs := []struct {
		name       string
		input      affix.SuffixDto
		categories resolver.Resolver[affix.SuffixCategory]
		expected   affix.Suffix
		errs       []error
	}{
		{name: "round trip", input: valid.ToDto(), categories: categories, expected: valid},
		{
			name:       "invalid dto",
			input:      affix.SuffixDto{Affix: affix.Affix{Id: id.New()}},
			categories: categories,
			errs:       []error{affix.ErrInvalid, affix.ErrName, affix.ErrCategory, id.ErrInvalid},
		},
		{
			name:       "unknown category",
			input:      affix.SuffixDto{Affix: valid.Affix, Category: id.New()},
			categories: categories,
			errs:       []error{resolver.ErrResolve, errNotFound},
		},
		{
			name:       "nil resolver",
			input:      valid.ToDto(),
			categories: nil,
			errs:       []error{resolver.ErrResolve, resolver.ErrNoResolver},
		},
		{
			name:       "resolved category is invalid",
			input:      affix.SuffixDto{Affix: valid.Affix, Category: invalidCategory.Id},
			categories: categories,
			errs:       []error{affix.ErrInvalid, affix.ErrCategory, affix.ErrInvalidCategory},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := tc.input.ToModel(tc.categories)

			if len(tc.errs) > 0 {
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, affix.Suffix{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_SuffixDto_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawCategoryId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"

	tcs := []struct {
		name     string
		input    affix.SuffixDto
		expected string
	}{
		{
			name:     "absent description is omitted",
			input:    affix.SuffixDto{Id: id.MustParse(rawId), Name: "Jr.", Category: id.MustParse(rawCategoryId)},
			expected: `{"id":"` + rawId + `","name":"Jr.","category":"` + rawCategoryId + `"}`,
		},
		{
			name: "all fields",
			input: affix.SuffixDto{
				Id:          id.MustParse(rawId),
				Name:        "Jr.",
				Description: optioner.Some("Junior"),
				Category:    id.MustParse(rawCategoryId),
			},
			expected: `{"id":"` + rawId + `","name":"Jr.","description":"Junior","category":"` + rawCategoryId + `"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual affix.SuffixDto
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_SuffixDto_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"Jr.","category":` + rawId + `}`, err: affix.ErrId},
		{name: "missing name", input: `{"id":` + rawId + `,"category":` + rawId + `}`, err: affix.ErrName},
		{name: "missing category", input: `{"id":` + rawId + `,"name":"Jr."}`, err: affix.ErrCategory},
		{name: "bad category id", input: `{"id":` + rawId + `,"name":"Jr.","category":"nope"}`, err: id.ErrInvalid},
		{name: "embedded category object", input: `{"id":` + rawId + `,"name":"Jr.","category":{"id":` + rawId + `}}`, err: id.ErrInvalid},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual affix.SuffixDto
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, affix.SuffixDto{}, actual)
		})
	}
}

func Test_SuffixDto_UnmarshalJSON_null(t *testing.T) {
	expected := affix.MustNewSuffix("Jr.", "", affix.MustNewSuffixCategory("Generational", "")).ToDto()
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
