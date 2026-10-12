// Package id implements the ID data type.
//
// Source of truth: domain-model/id/id.md
package id

import (
	"encoding/json"
	"uuid"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalid is returned when a value is not a valid UUID.
var ErrInvalid = errorer.New("invalid id")
var errInvalidFn = errorer.Func(ErrInvalid)

// Zero is the zero value (the Nil UUID). It is never returned by New.
var Zero = ID{}

// ID is a unique identifier, formatted as a version 4 UUID (RFC 9562).
//
// It serializes as a string, e.g. "f47ac10b-58cc-4372-a567-0e02b2c3d479".
type ID uuid.UUID

// String returns the lowercase hex-and-dash representation of the ID.
func (this ID) String() string {
	return uuid.UUID(this).String()
}

// IsZero reports whether the ID is the zero value.
func (this ID) IsZero() bool {
	return this == Zero
}

// Validate returns an error wrapping ErrInvalid if the ID is the zero value.
// Any non-zero ID is a valid UUID, since every constructor goes through Parse or New.
func (this ID) Validate() error {
	if this.IsZero() {
		return errInvalidFn("must not be zero")
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the ID serializes as a string.
func (this ID) MarshalText() ([]byte, error) {
	return uuid.UUID(this).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *ID) UnmarshalText(b []byte) error {
	id, err := Parse(string(b))

	if err != nil {
		return err
	}

	*this = id
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the ID as a JSON string.
func (this ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.String())
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a valid UUID. A JSON null leaves the ID unchanged.
func (this *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// New returns a new random (version 4) ID.
func New() ID {
	return ID(uuid.NewV4())
}

// Parse returns the ID represented by s.
// It returns an error wrapping ErrInvalid if s is not a valid UUID.
func Parse(s string) (ID, error) {
	u, err := uuid.Parse(s)

	if err != nil {
		return Zero, errInvalidFn("%q: %w", s, err)
	}

	return ID(u), nil
}

// MustParse is like Parse but panics if s is not a valid UUID.
func MustParse(s string) ID {
	id, err := Parse(s)

	if err != nil {
		panic(err)
	}

	return id
}
