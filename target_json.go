package thunder

type JsonTarget struct {
	ID     int    `json:"id"`
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

func SerializeTarget(s *Target) (*JsonTarget, error) {
	return &JsonTarget{
		ID:     s.ID,
		Driver: s.Driver.DriverInfo().ID,
		Config: s.Config.(any),
	}, nil
}

func UnserializeTarget(js *JsonTarget) (*Target, error) {
	driver, err := GetTargetDriver(js.Driver)
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

	return &Target{
		ID:     js.ID,
		Driver: driverInstance,
		Config: typedConfig,
	}, nil
}
