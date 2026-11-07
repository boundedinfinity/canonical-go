package fraction

import (
	"fmt"
	"strconv"

	"github.com/boundedinfinity/go-commoner/idiomatic/mather"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

// ----------------------------------------------------------------------------------------------------
// Constructors
// ----------------------------------------------------------------------------------------------------

func New[T mather.Integer](numerator, denominator T) Fraction[T] {
	return Fraction[T]{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

// FromFloat creates a Fraction from a floating point number.
func FromFloat[T mather.Integer, F mather.Float](n F) Fraction[T] {
	denominator := mather.FractionalPartPlace(n)
	numerator := n * F(denominator)

	return Fraction[T]{
		Numerator:   T(numerator),
		Denominator: T(denominator),
	}
}

// ----------------------------------------------------------------------------------------------------
// Type
// ----------------------------------------------------------------------------------------------------

type Fraction[T mather.Integer] struct {
	Numerator              T      `json:"numerator"`
	Denominator            T      `json:"denominator"`
	NumeratorDescription   string `json:"numerator-description"`
	DenominatorDescription string `json:"denominator-description"`
}

func (t Fraction[T]) String() string {
	return fmt.Sprintf("%d/%d", t.Numerator, t.Denominator)
}

func (t Fraction[T]) Description() string {
	if t.NumeratorDescription == "" || t.DenominatorDescription == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", t.NumeratorDescription, t.DenominatorDescription)
}

func (t Fraction[T]) Float() float64 {
	return float64(t.Numerator / t.Denominator)
}

func (t Fraction[T]) Copy() Fraction[T] {
	return Fraction[T]{
		Numerator:              t.Numerator,
		Denominator:            t.Denominator,
		NumeratorDescription:   t.NumeratorDescription,
		DenominatorDescription: t.DenominatorDescription,
	}
}

func (t Fraction[T]) Reduce() Fraction[T] {
	gcd := mather.GreatestCommonFactor(t.Numerator, t.Denominator)

	return Fraction[T]{
		Numerator:              t.Numerator / gcd,
		Denominator:            t.Denominator / gcd,
		NumeratorDescription:   t.NumeratorDescription,
		DenominatorDescription: t.DenominatorDescription,
	}
}

func (t Fraction[T]) IsZero() bool {
	return IsZero(t)
}

func (t Fraction[T]) Enumerate(l, h T) []Fraction[T] {
	var items []Fraction[T]
	item := t.Reduce()
	l = mather.Max(l, item.Denominator)

	for i := l; i <= h; i <<= 1 {
		items = append(items, item)
		item = New(item.Numerator*2, item.Denominator*2)
	}

	return items
}

// ----------------------------------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------------------------------

func IsZero[T mather.Integer](elem Fraction[T]) bool {
	var zero Fraction[T]
	return elem == zero
}

func Component[T mather.Float](n T) int {
	s := fmt.Sprintf("%v", n)
	comps := stringer.Split(s, ".")
	var i int

	if len(comps) > 1 {
		s = comps[1]
		x, err := strconv.Atoi(s)

		if err != nil {
			panic(err)
		}

		i = x
	} else {
		i = 0
	}

	return i
}

func Magnitude[T mather.Float](n T) int {
	s := fmt.Sprintf("%v", n)
	comps := stringer.Split(s, ".")
	var size int

	if len(comps) > 1 {
		s = stringer.Split(s, ".")[1]
		size = len(s)
	} else {
		s = "0"
		size = 0
	}

	return size
}
