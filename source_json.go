package thunder

type JsonSource struct {
	ID     int    `json:"id"`
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

func SerializeSource(s *Source) (*JsonSource, error) {
	return &JsonSource{
		ID:     s.ID,
		Driver: s.Driver.DriverInfo().ID,
		Config: s.Config.(any),
	}, nil
}

func UnserializeSource(js *JsonSource) (*Source, error) {
	driver, err := GetSourceDriver(js.Driver)
	if err != nil {
		return nil, err
	}

	typedConfig, err := ConvertToDynamicConfig(&driver.Config, js.Config)
	if err != nil {
		return nil, err
	}

	driverInstance, err := driver.New(typedConfig)
	if err != nil {
		return nil, err
	}

	return &Source{
		ID:     js.ID,
		Driver: driverInstance,
		Config: typedConfig,
	}, nil
}
