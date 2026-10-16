package affix_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/stretchr/testify/assert"
)

func Test_NewSuffixCategory(t *testing.T) {
	actual, err := affix.NewSuffixCategory("Generational", "")

	assert.NoError(t, err)
	assert.False(t, actual.Id.IsZero())
	assert.Equal(t, "Generational", actual.Name)
	assert.True(t, actual.Description.Empty())
}

func Test_NewSuffixCategory_invalid(t *testing.T) {
	actual, err := affix.NewSuffixCategory("", "a")

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrName)
	assert.ErrorIs(t, err, affix.ErrDescription)
	assert.Equal(t, affix.SuffixCategory{}, actual)
}

func Test_SuffixCategory_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	input := affix.SuffixCategory{Id: id.MustParse(rawId), Name: "Generational"}
	expected := `{"id":"` + rawId + `","name":"Generational"}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual affix.SuffixCategory
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_SuffixCategory_UnmarshalJSON_invalid(t *testing.T) {
	var actual affix.SuffixCategory
	err := json.Unmarshal([]byte(`{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":""}`), &actual)

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrName)
	assert.Equal(t, affix.SuffixCategory{}, actual)
}
