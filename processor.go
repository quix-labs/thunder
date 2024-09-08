package thunder

type Processor struct {
	Source int    `json:"source"`
	Table  string `json:"table"`

	Mapping Mapping `json:"mapping"`

	Index string `json:"index"`
}
