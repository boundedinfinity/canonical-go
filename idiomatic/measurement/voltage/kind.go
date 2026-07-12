package voltage

type Kind string

func (s Kind) String() string {
	return string(s)
}

var Kinds = kinds{
	AlternatingCurrent: "measurement.kind.voltage.alternating-current",
	DirectCurrent:      "measurement.kind.voltage.direct-current",
}

type kinds struct {
	AlternatingCurrent Kind
	DirectCurrent      Kind
}
