package voltage

import "fmt"

type Voltage struct {
	Kind  Kind    `json:"kind"`
	Value float64 `json:"value"`
}

func (this Voltage) String() string {
	var suffix string

	switch this.Kind {
	case Kinds.AlternatingCurrent:
		suffix = "AC"
	case Kinds.DirectCurrent:
		suffix = "DC"
	}

	return fmt.Sprintf("%f V%s", this.Value, suffix)
}

func AlternatingCurrent(value float64) Voltage {
	return Voltage{
		Kind:  Kinds.AlternatingCurrent,
		Value: value,
	}
}

func AC(value float64) Voltage {
	return AlternatingCurrent(value)
}

func DirectCurrent(value float64) Voltage {
	return Voltage{
		Kind:  Kinds.DirectCurrent,
		Value: value,
	}
}

func DC(value float64) Voltage {
	return DirectCurrent(value)
}
