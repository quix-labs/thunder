package thunder

import "github.com/quix-labs/thunder/utils"

type JsonTarget struct {
	ID     string `json:"id"`
	Driver string `json:"driver"`
	Config any    `json:"config"`
}

func SerializeTarget(t *Target) (*JsonTarget, error) {
	return &JsonTarget{
		ID:     t.ID,
		Driver: t.Driver.ID(),
		Config: t.Config.(any),
	}, nil
}

func UnserializeTarget(js *JsonTarget) (*Target, error) {
	driver, err := TargetDrivers.Get(js.Driver)
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

	return &Target{
		ID:     js.ID,
		Driver: driverInstance,
		Config: typedConfig,
	}, nil
}
