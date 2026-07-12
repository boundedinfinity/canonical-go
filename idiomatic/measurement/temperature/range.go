package temperature

import "fmt"

type Range struct {
	Min Temperature `json:"min"`
	Max Temperature `json:"max"`
}

func (this Range) String() string {
	return fmt.Sprintf("%s to %s", this.Min, this.Max)
}

func (this Range) In(temp Temperature) bool {
	r := this

	if temp.Scale != this.Min.Scale || temp.Scale != this.Max.Scale {
		r = this.ToScale(temp.Scale)
	}

	return temp.Value >= r.Min.Value && temp.Value <= r.Max.Value
}

func (this Range) ToScale(scale Scale) Range {
	return Range{
		Min: this.Min.ToScale(scale),
		Max: this.Max.ToScale(scale),
	}
}

func (this Range) Celsius() Range {
	return this.ToScale(Scales.Celsius)
}

func (this Range) Fahrenheit() Range {
	return this.ToScale(Scales.Fahrenheit)
}

func (this Range) Kelvin() Range {
	return this.ToScale(Scales.Kelvin)
}

func (this Range) Rankine() Range {
	return this.ToScale(Scales.Rankine)
}

func (this Range) Delisle() Range {
	return this.ToScale(Scales.Delisle)
}

func (this Range) Newton() Range {
	return this.ToScale(Scales.Newton)
}

func (this Range) Reaumur() Range {
	return this.ToScale(Scales.Reaumur)
}

func (this Range) Romer() Range {
	return this.ToScale(Scales.Romer)
}
