package label_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/label"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_LabelDto_round_trip(t *testing.T) {
	expected := label.MustNew("Date of Birth", []string{"DOB", "BD"}, "When a person was born")

	dto := expected.ToDto()
	assert.Equal(t, expected.Id, dto.Id)
	assert.Equal(t, expected.Name, dto.Name)
	assert.Equal(t, expected.Acronyms, dto.Acronyms)
	assert.Equal(t, expected.Description, dto.Description)

	actual, err := dto.ToModel()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_LabelDto_Validate(t *testing.T) {
	validId := id.MustParse("b12e7628-35b0-4ae5-b2d9-0f34f4828229")

	tcs := []struct {
		name  string
		input label.LabelDto
		errs  []error
	}{
		{name: "valid", input: label.LabelDto{Id: validId, Name: "n"}},
		{name: "zero id", input: label.LabelDto{Name: "n"}, errs: []error{label.ErrId}},
		{name: "name empty", input: label.LabelDto{Id: validId}, errs: []error{label.ErrName}},
		{name: "short acronym", input: label.LabelDto{Id: validId, Name: "n", Acronyms: optioner.Some([]string{"a"})}, errs: []error{label.ErrAcronyms}},
		{name: "empty acronyms", input: label.LabelDto{Id: validId, Name: "n", Acronyms: optioner.Some([]string{})}, errs: []error{label.ErrAcronyms}},
		{name: "empty description", input: label.LabelDto{Id: validId, Name: "n", Description: optioner.Some("")}, errs: []error{label.ErrDescription}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()
			_, modelErr := tc.input.ToModel()

			if len(tc.errs) > 0 {
				for _, e := range append(tc.errs, label.ErrInvalid) {
					assert.ErrorIs(tt, err, e)
					assert.ErrorIs(tt, modelErr, e)
				}
				return
			}

			assert.NoError(tt, err)
			assert.NoError(tt, modelErr)
		})
	}
}

func Test_LabelDto_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	input := label.LabelDto{Id: id.MustParse(rawId), Name: "Date of Birth", Acronyms: optioner.Some([]string{"DOB"})}
	expected := `{"id":"` + rawId + `","name":"Date of Birth","acronyms":["DOB"]}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual label.LabelDto
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_LabelDto_UnmarshalJSON_invalid(t *testing.T) {
	var actual label.LabelDto
	err := json.Unmarshal([]byte(`{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229"}`), &actual)

	assert.ErrorIs(t, err, label.ErrName)
	assert.Equal(t, label.LabelDto{}, actual)
}

func Test_LabelDto_UnmarshalJSON_null(t *testing.T) {
	expected := label.MustNew("n", nil, "").ToDto()
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
