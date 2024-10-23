package thunder

import "github.com/quix-labs/thunder/utils"

type JsonSource struct {
	ID     string `json:"id"`
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

func SerializeSource(s *Source) (*JsonSource, error) {
	return &JsonSource{
		ID:     s.ID,
		Driver: s.Driver.ID(),
		Config: s.Config.(any),
	}, nil
}

func UnserializeSource(js *JsonSource) (*Source, error) {
	driver, err := SourceDrivers.Get(js.Driver)
	if err != nil {
		return nil, err
	}

	config := driver.Config().Config
	typedConfig, err := utils.ConvertToDynamicConfig(config, js.Config)
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
