package affix_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/stretchr/testify/assert"
)

func Test_NewPrefixCategory(t *testing.T) {
	actual, err := affix.NewPrefixCategory("Honorific", "")

	assert.NoError(t, err)
	assert.False(t, actual.Id.IsZero())
	assert.Equal(t, "Honorific", actual.Name)
	assert.True(t, actual.Description.Empty())
}

func Test_NewPrefixCategory_invalid(t *testing.T) {
	actual, err := affix.NewPrefixCategory("", "a")

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrName)
	assert.ErrorIs(t, err, affix.ErrDescription)
	assert.Equal(t, affix.PrefixCategory{}, actual)
}

func Test_PrefixCategory_JSON(t *testing.T) {
	const rawId = "b12e7628-35b0-4ae5-b2d9-0f34f4828229"

	input := affix.PrefixCategory{Id: id.MustParse(rawId), Name: "Honorific"}
	expected := `{"id":"` + rawId + `","name":"Honorific"}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual affix.PrefixCategory
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_PrefixCategory_UnmarshalJSON_invalid(t *testing.T) {
	var actual affix.PrefixCategory
	err := json.Unmarshal([]byte(`{"id":"b12e7628-35b0-4ae5-b2d9-0f34f4828229","name":""}`), &actual)

	assert.ErrorIs(t, err, affix.ErrInvalidCategory)
	assert.ErrorIs(t, err, affix.ErrName)
	assert.Equal(t, affix.PrefixCategory{}, actual)
}
