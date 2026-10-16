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

func Test_NewCategory(t *testing.T) {
	tcs := []struct {
		name        string
		cname       string
		description string
		errs        []error
	}{
		{name: "name only", cname: "Honorific"},
		{name: "all fields", cname: "Honorific", description: "Titles of respect"},
		{name: "name min", cname: "a"},
		{name: "name max", cname: strings.Repeat("a", 30)},
		{name: "name max counts runes", cname: strings.Repeat("é", 30)},
		{name: "description min", cname: "a", description: "ab"},
		{name: "description max", cname: "a", description: strings.Repeat("a", 500)},
		{name: "description max counts runes", cname: "a", description: strings.Repeat("é", 500)},
		{name: "name empty", cname: "", errs: []error{affix.ErrName}},
		{name: "name too long", cname: strings.Repeat("a", 31), errs: []error{affix.ErrName}},
		{name: "description too short", cname: "a", description: "a", errs: []error{affix.ErrDescription}},
		{name: "description too long", cname: "a", description: strings.Repeat("a", 501), errs: []error{affix.ErrDescription}},
		{
			name:        "reports every failure",
			cname:       "",
			description: "a",
			errs:        []error{affix.ErrName, affix.ErrDescription},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := affix.NewCategory(tc.cname, tc.description)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, affix.ErrInvalidCategory)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, affix.Category{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.False(tt, actual.Id.IsZero())
			assert.Equal(tt, tc.cname, actual.Name)
			assert.Equal(tt, optioner.OfZero(tc.description), actual.Description)
		})
	}
}

func Test_Category_Validate_zero_id(t *testing.T) {
	err := affix.Category{Name: "a"}.Validate()

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrId)
	assert.ErrorIs(t, err, id.ErrInvalid)
}

func Test_Category_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	tcs := []struct {
		name     string
		input    affix.Category
		expected string
	}{
		{
			name:     "absent description is omitted",
			input:    affix.Category{Id: id.MustParse(rawId), Name: "Honorific"},
			expected: `{"id":"` + rawId + `","name":"Honorific"}`,
		},
		{
			name:     "all fields",
			input:    affix.Category{Id: id.MustParse(rawId), Name: "Honorific", Description: optioner.Some("Titles of respect")},
			expected: `{"id":"` + rawId + `","name":"Honorific","description":"Titles of respect"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual affix.Category
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_Category_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"a"}`, err: affix.ErrId},
		{name: "bad id", input: `{"id":"nope","name":"a"}`, err: id.ErrInvalid},
		{name: "missing name", input: `{"id":` + rawId + `}`, err: affix.ErrName},
		{name: "empty description", input: `{"id":` + rawId + `,"name":"a","description":""}`, err: affix.ErrDescription},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual affix.Category
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, affix.Category{}, actual)
		})
	}
}

func Test_Category_UnmarshalJSON_null(t *testing.T) {
	expected := affix.MustNewCategory("Honorific", "")
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
