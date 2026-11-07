package time

// Source of truth: domain-model/measurement/time/time-zone.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"
)

// TimeZoneDto is the data transfer object for TimeZone. TimeZone holds no
// other Entities, so TimeZoneDto has the same fields.
type TimeZoneDto struct {
	Identifier   string `json:"identifier"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	UtcOffset    string `json:"utc-offset"`
}

// Validate applies the TimeZone rules; see TimeZone.Validate.
func (this TimeZoneDto) Validate() error {
	return TimeZone(this).Validate()
}

// MarshalJSON implements json.Marshaler.
func (this TimeZoneDto) MarshalJSON() ([]byte, error) {
	return json.Marshal(TimeZone(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// TimeZoneDto, with the same rules as TimeZone.UnmarshalJSON.
// A JSON null leaves the TimeZoneDto unchanged.
func (this *TimeZoneDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var timeZone TimeZone

	if err := json.Unmarshal(data, &timeZone); err != nil {
		return err
	}

	*this = TimeZoneDto(timeZone)
	return nil
}

// ToModel returns the validated TimeZone for this TimeZoneDto.
func (this TimeZoneDto) ToModel() (TimeZone, error) {
	timeZone := TimeZone(this)

	if err := timeZone.Validate(); err != nil {
		return TimeZone{}, err
	}

	return timeZone, nil
}

// ToDto returns the TimeZoneDto for this TimeZone.
func (this TimeZone) ToDto() TimeZoneDto {
	return TimeZoneDto(this)
}
