package thunder

type JsonProcessor struct {
	ID int `json:"ID"`

	Source  int   `json:"source"`  // as source_id
	Targets []int `json:"targets"` // as targets_id

	Table       string      `json:"table"`
	PrimaryKeys []string    `json:"primary_keys"`
	Conditions  []Condition `json:"conditions"`

	Mapping JsonMapping `json:"mapping"`

	Index string `json:"index"`

	Enabled bool `json:"enabled"`
}

func SerializeProcessor(s *Processor) (*JsonProcessor, error) {
	jsonMapping, err := SerializeMapping(&s.Mapping)
	if err != nil {
		return nil, err
	}

	jp := JsonProcessor{
		ID:          s.ID,
		Source:      s.Source.ID,
		Targets:     make([]int, len(s.Targets)),
		Table:       s.Table,
		PrimaryKeys: s.PrimaryKeys,
		Conditions:  s.Conditions,
		Mapping:     *jsonMapping,
		Index:       s.Index,
		Enabled:     s.Enabled,
	}

	for i, target := range s.Targets {
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
	source, err := GetSource(jp.Source)
	if err != nil {
		return nil, err
	}
	p.Source = source

	// Load Targets
	p.Targets = make([]*Target, len(jp.Targets))
	for idx, targetId := range jp.Targets {
		target, err := GetTarget(targetId)
		if err != nil {
			return nil, err
		}
		p.Targets[idx] = target
	}

	return &p, nil
}
