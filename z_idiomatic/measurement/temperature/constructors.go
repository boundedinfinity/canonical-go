package temperature

func Celsius(value float64) Temperature {
	return Temperature{
		Scale: Scales.Celsius,
		Value: value,
	}
}

func Fahrenheit(value float64) Temperature {
	return Temperature{
		Scale: Scales.Fahrenheit,
		Value: value,
	}
}

func Kelvin(value float64) Temperature {
	return Temperature{
		Scale: Scales.Kelvin,
		Value: value,
	}
}

func Rankine(value float64) Temperature {
	return Temperature{
		Scale: Scales.Rankine,
		Value: value,
	}
}

func Delisle(value float64) Temperature {
	return Temperature{
		Scale: Scales.Delisle,
		Value: value,
	}
}

func Newton(value float64) Temperature {
	return Temperature{
		Scale: Scales.Newton,
		Value: value,
	}
}

func Reaumur(value float64) Temperature {
	return Temperature{
		Scale: Scales.Reaumur,
		Value: value,
	}
}

func Romer(value float64) Temperature {
	return Temperature{
		Scale: Scales.Romer,
		Value: value,
	}
}
