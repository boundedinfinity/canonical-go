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

func Test_NewLabelGroup(t *testing.T) {
	dob := label.MustNew("Date of Birth", []string{"DOB"}, "")
	ssn := label.MustNew("Social Security", []string{"SSN"}, "")

	tcs := []struct {
		name        string
		gname       string
		acronym     string
		description string
		labels      []label.Label
		errs        []error
	}{
		{name: "name only", gname: "Personal"},
		{name: "all fields", gname: "Personal", acronym: "PI", description: "Personal information", labels: []label.Label{dob, ssn}},
		{name: "name min", gname: "a"},
		{name: "name max", gname: strings.Repeat("a", 20)},
		{name: "name max counts runes", gname: strings.Repeat("é", 20)},
		{name: "name empty", gname: "", errs: []error{label.ErrName}},
		{name: "name too long", gname: strings.Repeat("a", 21), errs: []error{label.ErrName}},
		{name: "acronym min", gname: "n", acronym: "ab"},
		{name: "acronym max", gname: "n", acronym: strings.Repeat("a", 10)},
		{name: "acronym too short", gname: "n", acronym: "a", errs: []error{label.ErrAcronym}},
		{name: "acronym too long", gname: "n", acronym: strings.Repeat("a", 11), errs: []error{label.ErrAcronym}},
		{name: "description min", gname: "n", description: "a"},
		{name: "description max", gname: "n", description: strings.Repeat("a", 250)},
		{name: "description too long", gname: "n", description: strings.Repeat("a", 251), errs: []error{label.ErrDescription}},
		{name: "empty labels is absent", gname: "n", labels: []label.Label{}},
		{name: "invalid label", gname: "n", labels: []label.Label{dob, {Id: id.New()}}, errs: []error{label.ErrLabels, label.ErrInvalid, label.ErrName}},
		{
			name:        "reports every failure",
			gname:       "",
			acronym:     "a",
			description: strings.Repeat("a", 251),
			labels:      []label.Label{{}},
			errs:        []error{label.ErrName, label.ErrAcronym, label.ErrDescription, label.ErrLabels, label.ErrId},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := label.NewLabelGroup(tc.gname, tc.acronym, tc.description, tc.labels)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, label.ErrInvalidLabelGroup)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, label.LabelGroup{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.False(tt, actual.Id.IsZero())
			assert.Equal(tt, tc.gname, actual.Name)
			assert.Equal(tt, optioner.OfZero(tc.acronym), actual.Acronym)
			assert.Equal(tt, optioner.OfZero(tc.description), actual.Description)
			assert.Equal(tt, optioner.OfSlice(tc.labels), actual.Labels)
		})
	}
}

func Test_LabelGroup_Validate(t *testing.T) {
	validId := id.MustParse("b12e7628-35b0-4ae5-b2d9-0f34f4828229")

	tcs := []struct {
		name  string
		input label.LabelGroup
		errs  []error
	}{
		{name: "zero id", input: label.LabelGroup{Name: "n"}, errs: []error{label.ErrId, id.ErrInvalid}},
		{name: "present but empty labels", input: label.LabelGroup{Id: validId, Name: "n", Labels: optioner.Some([]label.Label{})}, errs: []error{label.ErrLabels}},
		{name: "empty description", input: label.LabelGroup{Id: validId, Name: "n", Description: optioner.Some("")}, errs: []error{label.ErrDescription}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()

			assert.ErrorIs(tt, err, label.ErrInvalidLabelGroup)
			for _, e := range tc.errs {
				assert.ErrorIs(tt, err, e)
			}
		})
	}
}

func Test_LabelGroup_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawLabelId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"

	dob := label.Label{Id: id.MustParse(rawLabelId), Name: "Date of Birth", Acronyms: optioner.Some([]string{"DOB"})}

	tcs := []struct {
		name     string
		input    label.LabelGroup
		expected string
	}{
		{
			name:     "name only omits optional fields",
			input:    label.LabelGroup{Id: id.MustParse(rawId), Name: "Personal"},
			expected: `{"id":"` + rawId + `","name":"Personal"}`,
		},
		{
			name: "all fields",
			input: label.LabelGroup{
				Id:          id.MustParse(rawId),
				Name:        "Personal",
				Acronym:     optioner.Some("PI"),
				Description: optioner.Some("Personal information"),
				Labels:      optioner.Some([]label.Label{dob}),
			},
			expected: `{"id":"` + rawId + `","name":"Personal","acronym":"PI","description":"Personal information",` +
				`"labels":[{"id":"` + rawLabelId + `","name":"Date of Birth","acronyms":["DOB"]}]}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual label.LabelGroup
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_LabelGroup_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"n"}`, err: label.ErrId},
		{name: "bad id", input: `{"id":"nope","name":"n"}`, err: id.ErrInvalid},
		{name: "missing name", input: `{"id":` + rawId + `}`, err: label.ErrName},
		{name: "short acronym", input: `{"id":` + rawId + `,"name":"n","acronym":"a"}`, err: label.ErrAcronym},
		{name: "empty description", input: `{"id":` + rawId + `,"name":"n","description":""}`, err: label.ErrDescription},
		{name: "empty labels", input: `{"id":` + rawId + `,"name":"n","labels":[]}`, err: label.ErrLabels},
		{name: "invalid label", input: `{"id":` + rawId + `,"name":"n","labels":[{"id":` + rawId + `}]}`, err: label.ErrName},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual label.LabelGroup
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, label.LabelGroup{}, actual)
		})
	}
}

func Test_LabelGroup_UnmarshalJSON_null(t *testing.T) {
	expected := label.MustNewLabelGroup("n", "", "", nil)
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
