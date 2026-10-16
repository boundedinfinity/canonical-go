package affix_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_NewSuffix(t *testing.T) {
	category := affix.MustNewSuffixCategory("Generational", "")

	tcs := []struct {
		name        string
		aname       string
		description string
		category    affix.SuffixCategory
		errs        []error
	}{
		{name: "name only", aname: "Jr.", category: category},
		{name: "all fields", aname: "Jr.", description: "Junior", category: category},
		{name: "name empty", aname: "", category: category, errs: []error{affix.ErrName}},
		{name: "name too long", aname: strings.Repeat("a", 31), category: category, errs: []error{affix.ErrName}},
		{name: "description too short", aname: "Jr.", description: "a", category: category, errs: []error{affix.ErrDescription}},
		{name: "missing category", aname: "Jr.", errs: []error{affix.ErrCategory, affix.ErrInvalidCategory, affix.ErrName, id.ErrInvalid}},
		{
			name:     "invalid category",
			aname:    "Jr.",
			category: affix.SuffixCategory{Id: id.New(), Name: "Generational", Description: optioner.Some("a")},
			errs:     []error{affix.ErrCategory, affix.ErrInvalidCategory, affix.ErrDescription},
		},
		{
			name:        "reports every failure",
			aname:       "",
			description: "a",
			errs:        []error{affix.ErrName, affix.ErrDescription, affix.ErrCategory},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := affix.NewSuffix(tc.aname, tc.description, tc.category)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, affix.ErrInvalid)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, affix.Suffix{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.False(tt, actual.Id.IsZero())
			assert.Equal(tt, tc.aname, actual.Name)
			assert.Equal(tt, optioner.OfZero(tc.description), actual.Description)
			assert.Equal(tt, tc.category, actual.Category)
		})
	}
}

func Test_Suffix_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawCategoryId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"

	category := affix.SuffixCategory{Id: id.MustParse(rawCategoryId), Name: "Generational"}
	categoryJSON := `{"id":"` + rawCategoryId + `","name":"Generational"}`

	tcs := []struct {
		name     string
		input    affix.Suffix
		expected string
	}{
		{
			name:     "absent description is omitted",
			input:    affix.Suffix{Id: id.MustParse(rawId), Name: "Jr.", Category: category},
			expected: `{"id":"` + rawId + `","name":"Jr.","category":` + categoryJSON + `}`,
		},
		{
			name: "all fields",
			input: affix.Suffix{
				Id:          id.MustParse(rawId),
				Name:        "Jr.",
				Description: optioner.Some("Junior"),
				Category:    category,
			},
			expected: `{"id":"` + rawId + `","name":"Jr.","description":"Junior","category":` + categoryJSON + `}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual affix.Suffix
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_Suffix_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`
	const category = `{"id":` + rawId + `,"name":"Generational"}`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"Jr.","category":` + category + `}`, err: affix.ErrId},
		{name: "missing name", input: `{"id":` + rawId + `,"category":` + category + `}`, err: affix.ErrName},
		{name: "empty description", input: `{"id":` + rawId + `,"name":"Jr.","description":"","category":` + category + `}`, err: affix.ErrDescription},
		{name: "missing category", input: `{"id":` + rawId + `,"name":"Jr."}`, err: affix.ErrCategory},
		{name: "invalid category", input: `{"id":` + rawId + `,"name":"Jr.","category":{"id":` + rawId + `,"name":""}}`, err: affix.ErrInvalidCategory},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual affix.Suffix
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, affix.Suffix{}, actual)
		})
	}
}

func Test_Suffix_UnmarshalJSON_null(t *testing.T) {
	expected := affix.MustNewSuffix("Jr.", "", affix.MustNewSuffixCategory("Generational", ""))
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
