package label_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/label"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_LabelGroupDto_ToDto(t *testing.T) {
	dob := label.MustNew("Date of Birth", []string{"DOB"}, "")
	ssn := label.MustNew("Social Security", []string{"SSN"}, "")

	model := label.MustNewLabelGroup("Personal", "PI", "Personal information", []label.Label{dob, ssn})
	dto := model.ToDto()

	assert.Equal(t, model.Id, dto.Id)
	assert.Equal(t, model.Name, dto.Name)
	assert.Equal(t, model.Acronym, dto.Acronym)
	assert.Equal(t, model.Description, dto.Description)
	assert.Equal(t, optioner.Some([]id.ID{dob.Id, ssn.Id}), dto.Labels)

	assert.True(t, label.MustNewLabelGroup("n", "", "", nil).ToDto().Labels.Empty())
}

func Test_LabelGroupDto_ToModel(t *testing.T) {
	errNotFound := errors.New("not found")
	dob := label.MustNew("Date of Birth", []string{"DOB"}, "")
	badLabel := label.Label{Id: id.New()}

	labels := resolver.Func[label.Label](func(i id.ID) (label.Label, error) {
		switch i {
		case dob.Id:
			return dob, nil
		case badLabel.Id:
			return badLabel, nil
		}

		return label.Label{}, errNotFound
	})

	full := label.MustNewLabelGroup("Personal", "PI", "", []label.Label{dob})
	bare := label.MustNewLabelGroup("Personal", "", "", nil)

	tcs := []struct {
		name     string
		input    label.LabelGroupDto
		labels   resolver.Resolver[label.Label]
		expected label.LabelGroup
		errs     []error
	}{
		{name: "round trip", input: full.ToDto(), labels: labels, expected: full},
		{name: "no labels needs no resolver", input: bare.ToDto(), expected: bare},
		{
			name:  "invalid dto",
			input: label.LabelGroupDto{Id: id.New(), Labels: optioner.Some([]id.ID{})},
			errs:  []error{label.ErrInvalidLabelGroup, label.ErrName, label.ErrLabels},
		},
		{
			name:   "unknown label",
			input:  label.LabelGroupDto{Id: id.New(), Name: "n", Labels: optioner.Some([]id.ID{dob.Id, id.New()})},
			labels: labels,
			errs:   []error{resolver.ErrResolve, errNotFound},
		},
		{
			name:  "missing resolver",
			input: full.ToDto(),
			errs:  []error{resolver.ErrResolve, resolver.ErrNoResolver},
		},
		{
			name:   "resolved label is invalid",
			input:  label.LabelGroupDto{Id: id.New(), Name: "n", Labels: optioner.Some([]id.ID{badLabel.Id})},
			labels: labels,
			errs:   []error{label.ErrInvalidLabelGroup, label.ErrLabels, label.ErrInvalid, label.ErrName},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := tc.input.ToModel(tc.labels)

			if len(tc.errs) > 0 {
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, label.LabelGroup{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_LabelGroupDto_Validate(t *testing.T) {
	validId := id.MustParse("b12e7628-35b0-4ae5-b2d9-0f34f4828229")

	tcs := []struct {
		name  string
		input label.LabelGroupDto
		errs  []error
	}{
		{name: "name only", input: label.LabelGroupDto{Id: validId, Name: "n"}},
		{name: "with labels", input: label.LabelGroupDto{Id: validId, Name: "n", Labels: optioner.Some([]id.ID{validId})}},
		{name: "zero id", input: label.LabelGroupDto{Name: "n"}, errs: []error{label.ErrId}},
		{name: "name empty", input: label.LabelGroupDto{Id: validId}, errs: []error{label.ErrName}},
		{name: "short acronym", input: label.LabelGroupDto{Id: validId, Name: "n", Acronym: optioner.Some("a")}, errs: []error{label.ErrAcronym}},
		{name: "empty description", input: label.LabelGroupDto{Id: validId, Name: "n", Description: optioner.Some("")}, errs: []error{label.ErrDescription}},
		{name: "empty labels", input: label.LabelGroupDto{Id: validId, Name: "n", Labels: optioner.Some([]id.ID{})}, errs: []error{label.ErrLabels}},
		{name: "zero label id", input: label.LabelGroupDto{Id: validId, Name: "n", Labels: optioner.Some([]id.ID{validId, {}})}, errs: []error{label.ErrLabels, id.ErrInvalid}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, label.ErrInvalidLabelGroup)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				return
			}

			assert.NoError(tt, err)
		})
	}
}

func Test_LabelGroupDto_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"
	const rawLabelId = "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"

	tcs := []struct {
		name     string
		input    label.LabelGroupDto
		expected string
	}{
		{
			name:     "name only omits optional fields",
			input:    label.LabelGroupDto{Id: id.MustParse(rawId), Name: "Personal"},
			expected: `{"id":"` + rawId + `","name":"Personal"}`,
		},
		{
			name: "all fields",
			input: label.LabelGroupDto{
				Id:          id.MustParse(rawId),
				Name:        "Personal",
				Acronym:     optioner.Some("PI"),
				Description: optioner.Some("Personal information"),
				Labels:      optioner.Some([]id.ID{id.MustParse(rawLabelId)}),
			},
			expected: `{"id":"` + rawId + `","name":"Personal","acronym":"PI","description":"Personal information",` +
				`"labels":["` + rawLabelId + `"]}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual label.LabelGroupDto
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_LabelGroupDto_UnmarshalJSON_invalid(t *testing.T) {
	const rawId = `"b12e7628-35b0-4ae5-b2d9-0f34f4828229"`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing id", input: `{"name":"n"}`, err: label.ErrId},
		{name: "missing name", input: `{"id":` + rawId + `}`, err: label.ErrName},
		{name: "empty labels", input: `{"id":` + rawId + `,"name":"n","labels":[]}`, err: label.ErrLabels},
		{name: "bad label id", input: `{"id":` + rawId + `,"name":"n","labels":["nope"]}`, err: id.ErrInvalid},
		{name: "embedded label object", input: `{"id":` + rawId + `,"name":"n","labels":[{"id":` + rawId + `,"name":"n"}]}`, err: id.ErrInvalid},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual label.LabelGroupDto
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, label.LabelGroupDto{}, actual)
		})
	}
}

func Test_LabelGroupDto_UnmarshalJSON_null(t *testing.T) {
	expected := label.MustNewLabelGroup("n", "", "", nil).ToDto()
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
