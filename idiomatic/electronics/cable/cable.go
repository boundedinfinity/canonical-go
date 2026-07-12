package cable

type Cable struct {
	Wires []Wire `json:"wires"`
}

type Wire struct {
	Size int `json:"size"`
}
