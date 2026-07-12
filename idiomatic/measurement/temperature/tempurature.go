package temperature

import (
	"fmt"
)

type Temperature struct {
	Scale Scale   `json:"scale"`
	Value float64 `json:"value"`
}

func (this Temperature) String() string {
	return fmt.Sprintf("%f %s", this.Value, clean(this.Scale))
}

func (this Temperature) Celsius() Temperature {
	return this.ToScale(Scales.Celsius)
}

func (this Temperature) Fahrenheit() Temperature {
	return this.ToScale(Scales.Fahrenheit)
}

func (this Temperature) Kelvin() Temperature {
	return this.ToScale(Scales.Kelvin)
}

func (this Temperature) Rankine() Temperature {
	return this.ToScale(Scales.Rankine)
}

func (this Temperature) Delisle() Temperature {
	return this.ToScale(Scales.Delisle)
}

func (this Temperature) Newton() Temperature {
	return this.ToScale(Scales.Newton)
}

func (this Temperature) Reaumur() Temperature {
	return this.ToScale(Scales.Reaumur)
}

func (this Temperature) Romer() Temperature {
	return this.ToScale(Scales.Romer)
}

func (this Temperature) ToScale(scale Scale) Temperature {
	if this.Scale == scale {
		return Temperature{
			Scale: this.Scale,
			Value: this.Value,
		}
	}

	conversion := fmt.Sprintf("%s to %s", clean(this.Scale), clean(scale))

	if fn, ok := sm[conversion]; ok {
		return Temperature{
			Scale: scale,
			Value: fn(this.Value),
		}
	}

	panic("invalid scale conversion: " + conversion)
}

var sm = map[string]func(float64) float64{
	"celsius to fahrenheit": celsius2fahrenheit,
	"celsius to kelvin":     celsius2kelvin,
	"celsius to rankine":    celsius2rankine,
	"celsius to delisle":    celsius2delisle,
	"celsius to newton":     celsius2newton,
	"celsius to reaumur":    celsius2reaumur,
	"celsius to romer":      celsius2romer,

	"delisle to celsius":    delisle2celsius,
	"delisle to fahrenheit": delisle2fahrenheit,
	"delisle to kelvin":     delisle2kelvin,
	"delisle to rankine":    delisle2rankine,
	"delisle to newton":     delisle2newton,
	"delisle to reaumur":    delisle2reaumur,
	"delisle to romer":      delisle2romer,

	"fahrenheit to celsius": fahrenheit2celsius,
	"fahrenheit to kelvin":  fahrenheit2kelvin,
	"fahrenheit to rankine": fahrenheit2rankine,
	"fahrenheit to delisle": fahrenheit2delisle,
	"fahrenheit to newton":  fahrenheit2newton,
	"fahrenheit to reaumur": fahrenheit2reaumur,
	"fahrenheit to romer":   fahrenheit2romer,

	"kelvin to celsius":    kelvin2celsius,
	"kelvin to fahrenheit": kelvin2fahrenheit,
	"kelvin to rankine":    kelvin2rankine,
	"kelvin to delisle":    kelvin2delisle,
	"kelvin to newton":     kelvin2newton,
	"kelvin to reaumur":    kelvin2reaumur,
	"kelvin to romer":      kelvin2romer,

	"newton to celsius":    newton2celsius,
	"newton to fahrenheit": newton2fahrenheit,
	"newton to kelvin":     newton2kelvin,
	"newton to rankine":    newton2rankine,
	"newton to delisle":    newton2delisle,
	"newton to reaumur":    newton2reaumur,
	"newton to romer":      newton2romer,

	"rankine to celsius":    rankine2celsius,
	"rankine to fahrenheit": rankine2fahrenheit,
	"rankine to kelvin":     rankine2kelvin,
	"rankine to delisle":    rankine2delisle,
	"rankine to newton":     rankine2newton,
	"rankine to reaumur":    rankine2reaumur,
	"rankine to romer":      rankine2romer,

	"reaumur to celsius":    reaumur2celsius,
	"reaumur to fahrenheit": reaumur2fahrenheit,
	"reaumur to kelvin":     reaumur2kelvin,
	"reaumur to rankine":    reaumur2rankine,
	"reaumur to delisle":    reaumur2delisle,
	"reaumur to newton":     reaumur2newton,
	"reaumur to romer":      reaumur2romer,

	"romer to celsius":    romer2celsius,
	"romer to fahrenheit": romer2fahrenheit,
	"romer to kelvin":     romer2kelvin,
	"romer to rankine":    romer2rankine,
	"romer to delisle":    romer2delisle,
	"romer to newton":     romer2newton,
	"romer to reaumur":    romer2reaumur,
}

func celsius2fahrenheit(v float64) float64 { return (v * 9 / 5) + 32 }
func celsius2kelvin(v float64) float64     { return v + 273.15 }
func celsius2rankine(v float64) float64    { return (v + 273.15) * 9 / 5 }
func celsius2delisle(v float64) float64    { return (100 - v) * 3 / 2 }
func celsius2newton(v float64) float64     { return v * 33 / 100 }
func celsius2reaumur(v float64) float64    { return v * 4 / 5 }
func celsius2romer(v float64) float64      { return v*21/40 + 7.5 }

func delisle2celsius(v float64) float64    { return 100 - (v * 2 / 3) }
func delisle2fahrenheit(v float64) float64 { return 212 - (v * 6 / 5) }
func delisle2kelvin(v float64) float64     { return 373.15 - (v * 2 / 3) }
func delisle2rankine(v float64) float64    { return 671.67 - (v * 6 / 5) }
func delisle2newton(v float64) float64     { return 33 - (v * 11 / 50) }
func delisle2reaumur(v float64) float64    { return 100 - (v * 8 / 15) }
func delisle2romer(v float64) float64      { return 60 - (v * 7 / 20) }

func fahrenheit2celsius(v float64) float64 { return (v - 32) * 5 / 9 }
func fahrenheit2kelvin(v float64) float64  { return (v + 459.67) * 5 / 9 }
func fahrenheit2rankine(v float64) float64 { return v + 459.67 }
func fahrenheit2delisle(v float64) float64 { return (212 - v) * 5 / 6 }
func fahrenheit2newton(v float64) float64  { return (v - 32) * 11 / 60 }
func fahrenheit2reaumur(v float64) float64 { return (v - 32) * 4 / 9 }
func fahrenheit2romer(v float64) float64   { return (v-32)*7/24 + 7.5 }

func kelvin2celsius(v float64) float64    { return v - 273.15 }
func kelvin2fahrenheit(v float64) float64 { return (v * 9 / 5) - 459.67 }
func kelvin2rankine(v float64) float64    { return v * 9 / 5 }
func kelvin2delisle(v float64) float64    { return (373.15 - v) * 3 / 2 }
func kelvin2newton(v float64) float64     { return (v - 273.15) * 33 / 100 }
func kelvin2reaumur(v float64) float64    { return (v - 273.15) * 4 / 5 }
func kelvin2romer(v float64) float64      { return (v-273.15)*21/40 + 7.5 }

func rankine2celsius(v float64) float64    { return (v - 491.67) * 5 / 9 }
func rankine2fahrenheit(v float64) float64 { return v - 459.67 }
func rankine2kelvin(v float64) float64     { return v * 5 / 9 }
func rankine2delisle(v float64) float64    { return (671.67 - v) * 5 / 6 }
func rankine2newton(v float64) float64     { return (v - 491.67) * 11 / 60 }
func rankine2reaumur(v float64) float64    { return (v - 491.67) * 4 / 9 }
func rankine2romer(v float64) float64      { return (v-491.67)*7/24 + 7.5 }

func newton2celsius(v float64) float64    { return v * 100 / 33 }
func newton2fahrenheit(v float64) float64 { return (v * 60 / 11) + 32 }
func newton2kelvin(v float64) float64     { return (v * 100 / 33) + 273.15 }
func newton2rankine(v float64) float64    { return (v * 60 / 11) + 491.67 }
func newton2delisle(v float64) float64    { return (33 - v) * 50 / 11 }
func newton2reaumur(v float64) float64    { return v * 4 / 5 }
func newton2romer(v float64) float64      { return v*21/40 + 7.5 }

func reaumur2celsius(v float64) float64    { return v * 5 / 4 }
func reaumur2fahrenheit(v float64) float64 { return (v * 9 / 4) + 32 }
func reaumur2kelvin(v float64) float64     { return (v * 5 / 4) + 273.15 }
func reaumur2rankine(v float64) float64    { return (v * 9 / 4) + 491.67 }
func reaumur2delisle(v float64) float64    { return (80 - v) * 15 / 8 }
func reaumur2newton(v float64) float64     { return v * 33 / 40 }
func reaumur2romer(v float64) float64      { return v*21/32 + 7.5 }

func romer2celsius(v float64) float64    { return (v - 7.5) * 40 / 21 }
func romer2fahrenheit(v float64) float64 { return (v-7.5)*24/7 + 32 }
func romer2kelvin(v float64) float64     { return (v-7.5)*40/21 + 273.15 }
func romer2rankine(v float64) float64    { return (v-7.5)*24/7 + 491.67 }
func romer2delisle(v float64) float64    { return (60 - v) * 20 / 7 }
func romer2newton(v float64) float64     { return (v - 7.5) * 40 / 35 }
func romer2reaumur(v float64) float64    { return (v - 7.5) * 32 / 21 }
