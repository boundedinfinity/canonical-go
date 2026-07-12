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

func (this Range) Celcius() Range {
	return Range{
		Min: this.Min.Celsius(),
		Max: this.Max.Celsius(),
	}
}

func (this Range) Fahrenheit() Range {
	return Range{
		Min: this.Min.Fahrenheit(),
		Max: this.Max.Fahrenheit(),
	}
}

func (this Range) Kelvin() Range {
	return Range{
		Min: this.Min.Kelvin(),
		Max: this.Max.Kelvin(),
	}
}

func (this Range) Rankine() Range {
	return Range{
		Min: this.Min.Rankine(),
		Max: this.Max.Rankine(),
	}
}

func (this Range) Delisle() Range {
	return Range{
		Min: this.Min.Delisle(),
		Max: this.Max.Delisle(),
	}
}

func (this Range) Newton() Range {
	return Range{
		Min: this.Min.Newton(),
		Max: this.Max.Newton(),
	}
}

func (this Range) Reaumur() Range {
	return Range{
		Min: this.Min.Reaumur(),
		Max: this.Max.Reaumur(),
	}
}

func (this Range) Romer() Range {
	return Range{
		Min: this.Min.Romer(),
		Max: this.Max.Romer(),
	}
}
