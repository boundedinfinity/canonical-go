package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_TimeZoneDto_RoundTrip(t *testing.T) {
	model := validTimeZone()
	dto := model.ToDto()

	assert.Equal(t, model.Identifier, dto.Identifier)
	assert.NoError(t, dto.Validate())

	actual, err := dto.ToModel()
	assert.NoError(t, err)
	assert.Equal(t, model, actual)
}

func Test_TimeZoneDto_Validate(t *testing.T) {
	dto := validTimeZone().ToDto()
	dto.Abbreviation = "est"

	assert.ErrorIs(t, dto.Validate(), time.ErrAbbreviation)

	_, err := dto.ToModel()
	assert.ErrorIs(t, err, time.ErrInvalidTimeZone)
}

func Test_TimeZoneDto_JSON(t *testing.T) {
	expected := validTimeZone().ToDto()

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"identifier":"America/New_York","name":"Eastern Standard Time","abbreviation":"EST","utc-offset":"-05:00"}`, string(data))

	var actual time.TimeZoneDto
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"identifier":"America/New_York"}`), &actual), time.ErrName)
	assert.Equal(t, expected, actual, "invalid JSON leaves the DTO unchanged")

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual, "null leaves the DTO unchanged")
}
