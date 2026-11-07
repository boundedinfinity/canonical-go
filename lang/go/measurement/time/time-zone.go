package time

// Source of truth: domain-model/measurement/time/time-zone.md

import (
	"encoding/json"
	"regexp"
	"time"
	_ "time/tzdata"

	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	timeZoneIdentifierValidation   = validation.String[string]("identifier").Min(1).Max(100).Match(regexp.MustCompile(`^[A-Za-z0-9_+\-]+(/[A-Za-z0-9_+\-]+)*$`))
	timeZoneNameValidation         = validation.String[string]("name").Min(1).Max(100)
	timeZoneAbbreviationValidation = validation.String[string]("abbreviation").Match(regexp.MustCompile(`^[A-Z]{2,6}$`))
	timeZoneUtcOffsetValidation    = validation.String[string]("utc-offset").Match(regexp.MustCompile(`^[+-](0[0-9]|1[0-4]):(00|15|30|45)$`))

	// ErrInvalidTimeZone is wrapped by every TimeZone validation error.
	ErrInvalidTimeZone = errorer.New("invalid time zone")
	ErrIdentifier      = errorer.New("identifier")
	ErrName            = errorer.New("name")
	ErrAbbreviation    = errorer.New("abbreviation")
	ErrUtcOffset       = errorer.New("utc-offset")

	errTimeZoneIdentifierFn   = errorer.Func(ErrInvalidTimeZone, ErrIdentifier)
	errTimeZoneNameFn         = errorer.Func(ErrInvalidTimeZone, ErrName)
	errTimeZoneAbbreviationFn = errorer.Func(ErrInvalidTimeZone, ErrAbbreviation)
	errTimeZoneUtcOffsetFn    = errorer.Func(ErrInvalidTimeZone, ErrUtcOffset)
)

// TimeZone is a time zone identified by its IANA identifier, e.g.
// "America/New_York". Its natural key is Identifier; it has no Id.
//
// Name, Abbreviation, and UtcOffset describe the zone's standard time. Use
// Location for the offset at a given moment, including daylight saving time.
type TimeZone struct {
	Identifier   string `json:"identifier"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	UtcOffset    string `json:"utc-offset"`
}

// Location returns the time.Location for Identifier, which applies daylight
// saving time and historical changes. It returns an error wrapping
// ErrInvalidTimeZone and ErrIdentifier if the identifier is not in the IANA
// database.
func (this TimeZone) Location() (*time.Location, error) {
	location, err := time.LoadLocation(this.Identifier)

	if err != nil {
		return nil, errTimeZoneIdentifierFn("%q: %w", this.Identifier, err)
	}

	return location, nil
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidTimeZone and the field's sentinel
// (ErrIdentifier, etc.). Lengths are counted in characters (runes), not bytes.
//
// It checks the format of Identifier only; Location checks that it exists.
func (this TimeZone) Validate() error {
	var errs []error

	if err := timeZoneIdentifierValidation.Validate(this.Identifier); err != nil {
		errs = append(errs, errTimeZoneIdentifierFn("%w", err))
	}

	if err := timeZoneNameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errTimeZoneNameFn("%w", err))
	}

	if err := timeZoneAbbreviationValidation.Validate(this.Abbreviation); err != nil {
		errs = append(errs, errTimeZoneAbbreviationFn("%w", err))
	}

	if err := timeZoneUtcOffsetValidation.Validate(this.UtcOffset); err != nil {
		errs = append(errs, errTimeZoneUtcOffsetFn("%w", err))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this TimeZone) MarshalJSON() ([]byte, error) {
	type alias TimeZone
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// TimeZone. A JSON null leaves the TimeZone unchanged.
func (this *TimeZone) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias TimeZone
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	timeZone := TimeZone(decoded)

	if err := timeZone.Validate(); err != nil {
		return err
	}

	*this = timeZone
	return nil
}

// NewTimeZone returns a validated TimeZone.
func NewTimeZone(identifier, name, abbreviation, utcOffset string) (TimeZone, error) {
	timeZone := TimeZone{
		Identifier:   identifier,
		Name:         name,
		Abbreviation: abbreviation,
		UtcOffset:    utcOffset,
	}

	if err := timeZone.Validate(); err != nil {
		return TimeZone{}, err
	}

	return timeZone, nil
}

// MustNewTimeZone is like NewTimeZone but panics if the TimeZone is invalid.
func MustNewTimeZone(identifier, name, abbreviation, utcOffset string) TimeZone {
	timeZone, err := NewTimeZone(identifier, name, abbreviation, utcOffset)

	if err != nil {
		panic(err)
	}

	return timeZone
}
