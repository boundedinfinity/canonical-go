package time

// Source of truth: domain-model/measurement/time/quarter.md

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalidQuarter is returned when a value is not a valid quarter.
var ErrInvalidQuarter = errorer.New("invalid quarter")
var errInvalidQuarterFn = errorer.Func(ErrInvalidQuarter)

// Quarter is a calendar quarter of the year with no year, numbered from
// Q1 (1, January–March) to Q4 (4, October–December).
//
// It serializes as its name, e.g. "Q1".
type Quarter int

const (
	Q1 Quarter = iota + 1
	Q2
	Q3
	Q4
)

// ZeroQuarter is the zero value. It is not a valid Quarter.
var ZeroQuarter = Quarter(0)

// QuartersOfYear lists every Quarter in order, Q1 first.
var QuartersOfYear = []Quarter{Q1, Q2, Q3, Q4}

// Number returns the quarter number, 1 (Q1) through 4 (Q4).
func (this Quarter) Number() int {
	return int(this)
}

// Months returns the quarter's three months, in order.
// The result for an invalid Quarter is undefined.
func (this Quarter) Months() []Month {
	first := Month(this-1)*3 + 1
	return []Month{first, first + 1, first + 2}
}

// String returns the quarter's name, e.g. "Q1", or a placeholder such as
// "Quarter(0)" for an invalid value.
func (this Quarter) String() string {
	if this.Validate() != nil {
		return "Quarter(" + strconv.Itoa(int(this)) + ")"
	}

	return "Q" + strconv.Itoa(int(this))
}

// IsZero reports whether the Quarter is the zero value.
func (this Quarter) IsZero() bool {
	return this == ZeroQuarter
}

// Validate returns an error wrapping ErrInvalidQuarter if the Quarter is not
// between Q1 (1) and Q4 (4).
func (this Quarter) Validate() error {
	if this < Q1 || this > Q4 {
		return errInvalidQuarterFn("%d must be between 1 and 4", int(this))
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the Quarter serializes
// as its name. It returns the validation error for an invalid Quarter.
func (this Quarter) MarshalText() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *Quarter) UnmarshalText(b []byte) error {
	q, err := ParseQuarter(string(b))

	if err != nil {
		return err
	}

	*this = q
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the Quarter as a JSON
// string. It returns the validation error for an invalid Quarter.
func (this Quarter) MarshalJSON() ([]byte, error) {
	text, err := this.MarshalText()

	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a quarter name. A JSON null leaves the Quarter unchanged.
func (this *Quarter) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidQuarterFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// NewQuarter returns the quarter of t, in t's own location.
func NewQuarter(t time.Time) Quarter {
	return Quarter((int(t.Month())-1)/3 + 1)
}

// QuarterOfMonth returns the quarter that m falls in. It returns an error
// wrapping ErrInvalidMonth if m is invalid.
func QuarterOfMonth(m Month) (Quarter, error) {
	if err := m.Validate(); err != nil {
		return ZeroQuarter, err
	}

	return Quarter((int(m)-1)/3 + 1), nil
}

// ParseQuarter returns the Quarter named by s, ignoring case ("Q1", "q1").
// It returns an error wrapping ErrInvalidQuarter for anything else.
func ParseQuarter(s string) (Quarter, error) {
	for _, q := range QuartersOfYear {
		if strings.EqualFold(s, q.String()) {
			return q, nil
		}
	}

	return ZeroQuarter, errInvalidQuarterFn("%q", s)
}

// MustParseQuarter is like ParseQuarter but panics if s is not a quarter.
func MustParseQuarter(s string) Quarter {
	q, err := ParseQuarter(s)

	if err != nil {
		panic(err)
	}

	return q
}
