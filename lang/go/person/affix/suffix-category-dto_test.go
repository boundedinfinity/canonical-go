package affix_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/stretchr/testify/assert"
)

func Test_SuffixCategoryDto_round_trip(t *testing.T) {
	expected := affix.MustNewSuffixCategory("Generational", "A description")

	dto := expected.ToDto()
	assert.Equal(t, expected.Category, dto.Category)

	actual, err := dto.ToModel()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_SuffixCategoryDto_ToModel_invalid(t *testing.T) {
	actual, err := affix.SuffixCategoryDto{Name: "Generational"}.ToModel()

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrId)
	assert.Equal(t, affix.SuffixCategory{}, actual)
}

func Test_SuffixCategoryDto_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	input := affix.SuffixCategoryDto{Id: id.MustParse(rawId), Name: "Generational"}
	expected := `{"id":"` + rawId + `","name":"Generational"}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual affix.SuffixCategoryDto
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_SuffixCategoryDto_UnmarshalJSON_invalid(t *testing.T) {
	var actual affix.SuffixCategoryDto
	err := json.Unmarshal([]byte(`{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":""}`), &actual)

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrName)
	assert.Equal(t, affix.SuffixCategoryDto{}, actual)
}
