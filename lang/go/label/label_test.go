package label_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/label"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	tcs := []struct {
		name        string
		lname       string
		acronyms    []string
		description string
		errs        []error
	}{
		{name: "name only", lname: "Date of Birth"},
		{name: "all fields", lname: "Date of Birth", acronyms: []string{"DOB"}, description: "When a person was born"},
		{name: "multiple acronyms", lname: "Date of Birth", acronyms: []string{"DOB", "BD"}},
		{name: "empty acronyms is absent", lname: "n", acronyms: []string{}},
		{name: "name min", lname: "a"},
		{name: "name max", lname: strings.Repeat("a", 20)},
		{name: "name max counts runes", lname: strings.Repeat("é", 20)},
		{name: "name empty", lname: "", errs: []error{label.ErrName}},
		{name: "name too long", lname: strings.Repeat("a", 21), errs: []error{label.ErrName}},
		{name: "acronym min", lname: "n", acronyms: []string{"ab"}},
		{name: "acronym max", lname: "n", acronyms: []string{strings.Repeat("a", 10)}},
		{name: "acronym max counts runes", lname: "n", acronyms: []string{strings.Repeat("é", 10)}},
		{name: "acronym too short", lname: "n", acronyms: []string{"a"}, errs: []error{label.ErrAcronyms}},
		{name: "acronym too long", lname: "n", acronyms: []string{strings.Repeat("a", 11)}, errs: []error{label.ErrAcronyms}},
		{name: "bad acronym after good", lname: "n", acronyms: []string{"DOB", "a"}, errs: []error{label.ErrAcronyms}},
		{name: "description max", lname: "n", description: strings.Repeat("a", 250)},
		{name: "description too long", lname: "n", description: strings.Repeat("a", 251), errs: []error{label.ErrDescription}},
		{
			name:        "reports every failure",
			lname:       "",
			acronyms:    []string{"a"},
			description: strings.Repeat("a", 251),
			errs:        []error{label.ErrName, label.ErrAcronyms, label.ErrDescription},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := label.New(tc.lname, tc.acronyms, tc.description)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, label.ErrInvalid)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, label.Label{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.False(tt, actual.Id.IsZero())
			assert.Equal(tt, tc.lname, actual.Name)
			assert.Equal(tt, optioner.OfSlice(tc.acronyms), actual.Acronyms)
			assert.Equal(tt, optioner.OfZero(tc.description), actual.Description)
		})
	}
}

func Test_Validate_zero_id(t *testing.T) {
	err := label.Label{Name: "n"}.Validate()

	assert.ErrorIs(t, err, label.ErrInvalid)
	assert.ErrorIs(t, err, label.ErrId)
	assert.ErrorIs(t, err, id.ErrInvalid)
}

func Test_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	tcs := []struct {
		name     string
		input    label.Label
		expected string
	}{
		{
			name:     "name only omits optional fields",
			input:    label.Label{Id: id.MustParse(rawId), Name: "Date of Birth"},
			expected: `{"id":"` + rawId + `","name":"Date of Birth"}`,
		},
		{
			name: "all fields",
			input: label.Label{
				Id:          id.MustParse(rawId),
				Name:        "Date of Birth",
				Acronyms:    optioner.Some([]string{"DOB", "BD"}),
				Description: optioner.Some("When a person was born"),
			},
			expected: `{"id":"` + rawId + `","name":"Date of Birth","acronyms":["DOB","BD"],"description":"When a person was born"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual label.Label
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_UnmarshalJSON_invalid(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"n"}`, err: label.ErrId},
		{name: "bad id", input: `{"id":"nope","name":"n"}`, err: id.ErrInvalid},
		{name: "missing name", input: `{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229"}`, err: label.ErrName},
		{name: "short acronym", input: `{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":"n","acronyms":["a"]}`, err: label.ErrAcronyms},
		{name: "empty acronyms", input: `{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":"n","acronyms":[]}`, err: label.ErrAcronyms},
		{name: "empty description", input: `{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":"n","description":""}`, err: label.ErrDescription},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual label.Label
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, label.Label{}, actual)
		})
	}
}

func Test_UnmarshalJSON_null(t *testing.T) {
	expected := label.MustNew("n", nil, "")
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
