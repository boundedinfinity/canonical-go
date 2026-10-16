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

func Test_NewPrefix(t *testing.T) {
	category := affix.MustNewPrefixCategory("Honorific", "")

	tcs := []struct {
		name        string
		aname       string
		description string
		category    affix.PrefixCategory
		errs        []error
	}{
		{name: "name only", aname: "Dr.", category: category},
		{name: "all fields", aname: "Dr.", description: "Doctor", category: category},
		{name: "name empty", aname: "", category: category, errs: []error{affix.ErrName}},
		{name: "name too long", aname: strings.Repeat("a", 31), category: category, errs: []error{affix.ErrName}},
		{name: "description too short", aname: "Dr.", description: "a", category: category, errs: []error{affix.ErrDescription}},
		{name: "missing category", aname: "Dr.", errs: []error{affix.ErrCategory, affix.ErrInvalidCategory, affix.ErrName, id.ErrInvalid}},
		{
			name:     "invalid category",
			aname:    "Dr.",
			category: affix.PrefixCategory{Id: id.New(), Name: "Honorific", Description: optioner.Some("a")},
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
			actual, err := affix.NewPrefix(tc.aname, tc.description, tc.category)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, affix.ErrInvalid)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, affix.Prefix{}, actual)
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

func Test_Prefix_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawCategoryId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"

	category := affix.PrefixCategory{Id: id.MustParse(rawCategoryId), Name: "Honorific"}
	categoryJSON := `{"id":"` + rawCategoryId + `","name":"Honorific"}`

	tcs := []struct {
		name     string
		input    affix.Prefix
		expected string
	}{
		{
			name:     "absent description is omitted",
			input:    affix.Prefix{Id: id.MustParse(rawId), Name: "Dr.", Category: category},
			expected: `{"id":"` + rawId + `","name":"Dr.","category":` + categoryJSON + `}`,
		},
		{
			name: "all fields",
			input: affix.Prefix{
				Id:          id.MustParse(rawId),
				Name:        "Dr.",
				Description: optioner.Some("Doctor"),
				Category:    category,
			},
			expected: `{"id":"` + rawId + `","name":"Dr.","description":"Doctor","category":` + categoryJSON + `}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual affix.Prefix
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_Prefix_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`
	const category = `{"id":` + rawId + `,"name":"Honorific"}`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"Dr.","category":` + category + `}`, err: affix.ErrId},
		{name: "missing name", input: `{"id":` + rawId + `,"category":` + category + `}`, err: affix.ErrName},
		{name: "empty description", input: `{"id":` + rawId + `,"name":"Dr.","description":"","category":` + category + `}`, err: affix.ErrDescription},
		{name: "missing category", input: `{"id":` + rawId + `,"name":"Dr."}`, err: affix.ErrCategory},
		{name: "invalid category", input: `{"id":` + rawId + `,"name":"Dr.","category":{"id":` + rawId + `,"name":""}}`, err: affix.ErrInvalidCategory},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual affix.Prefix
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, affix.Prefix{}, actual)
		})
	}
}

func Test_Prefix_UnmarshalJSON_null(t *testing.T) {
	expected := affix.MustNewPrefix("Dr.", "", affix.MustNewPrefixCategory("Honorific", ""))
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
