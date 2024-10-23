package thunder

type JsonProcessor struct {
	ID string `json:"ID"`

	Source  string   `json:"source"`  // as source_id
	Targets []string `json:"targets"` // as targets_id

	Table       string      `json:"table"`
	PrimaryKeys []string    `json:"primary_keys"`
	Conditions  []Condition `json:"conditions"`

	Mapping JsonMapping `json:"mapping"`

	Index string `json:"index"`

	Enabled bool `json:"enabled"`
}

func SerializeProcessor(p *Processor) (*JsonProcessor, error) {
	jsonMapping, err := SerializeMapping(&p.Mapping)
	if err != nil {
		return nil, err
	}

	jp := JsonProcessor{
		ID:          p.ID,
		Source:      p.Source.ID,
		Targets:     make([]string, len(p.Targets)),
		Table:       p.Table,
		PrimaryKeys: p.PrimaryKeys,
		Conditions:  p.Conditions,
		Mapping:     *jsonMapping,
		Index:       p.Index,
		Enabled:     p.Enabled,
	}

	for i, target := range p.Targets {
		jp.Targets[i] = target.ID
	}

	return &jp, nil
}

func UnserializeProcessor(jp *JsonProcessor) (*Processor, error) {
	mapping, err := UnserializeMapping(&jp.Mapping, nil)
	if err != nil {
		return nil, err
	}
	p := Processor{
		ID:          jp.ID,
		Table:       jp.Table,
		PrimaryKeys: jp.PrimaryKeys,
		Conditions:  jp.Conditions,
		Mapping:     *mapping,
		Index:       jp.Index,
		Enabled:     jp.Enabled,
	}

	// Load source
	source, err := Sources.Get(jp.Source)
	if err != nil {
		return nil, err
	}
	p.Source = &source

	// Load Targets
	p.Targets = make([]*Target, len(jp.Targets))

	for idx, targetId := range jp.Targets {
		target, err := Targets.Get(targetId)
		if err != nil {
			return nil, err
		}
		p.Targets[idx] = &target
	}

	return &p, nil
}
