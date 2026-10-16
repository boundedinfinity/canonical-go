package affix_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/stretchr/testify/assert"
)

func Test_PrefixCategoryDto_round_trip(t *testing.T) {
	expected := affix.MustNewPrefixCategory("Honorific", "A description")

	dto := expected.ToDto()
	assert.Equal(t, expected.Category, dto.Category)

	actual, err := dto.ToModel()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_PrefixCategoryDto_ToModel_invalid(t *testing.T) {
	actual, err := affix.PrefixCategoryDto{Name: "Honorific"}.ToModel()

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrId)
	assert.Equal(t, affix.PrefixCategory{}, actual)
}

func Test_PrefixCategoryDto_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	input := affix.PrefixCategoryDto{Id: id.MustParse(rawId), Name: "Honorific"}
	expected := `{"id":"` + rawId + `","name":"Honorific"}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual affix.PrefixCategoryDto
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_PrefixCategoryDto_UnmarshalJSON_invalid(t *testing.T) {
	var actual affix.PrefixCategoryDto
	err := json.Unmarshal([]byte(`{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":""}`), &actual)

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrName)
	assert.Equal(t, affix.PrefixCategoryDto{}, actual)
}
