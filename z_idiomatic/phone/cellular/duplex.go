package cellular

type DuplexMode struct {
	Abbr string
	Name string
}

var DuplexModes = duplexModes{
	FDD: DuplexMode{
		Abbr: "FDD",
		Name: "Frequency Division Duplex",
	},
	TDD: DuplexMode{
		Abbr: "TDD",
		Name: "Time Division Duplex",
	},
	SDL: DuplexMode{
		Abbr: "FMSD",
		Name: "Frequency Modulation Supplimental Downlink",
	},
	SDO: DuplexMode{
		Abbr: "SDO",
		Name: "Standalone Downlink Only",
	},
}

type duplexModes struct {
	FDD DuplexMode
	TDD DuplexMode
	SDL DuplexMode
	SDO DuplexMode
}
