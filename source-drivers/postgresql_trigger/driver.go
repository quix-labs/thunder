package postgresql_trigger

import (
	_ "embed"
	"errors"
	"github.com/quix-labs/thunder"
)

func init() {
	thunder.RegisterSourceDriver(&Driver{})
}

type DriverConfig struct {
	//Host     string `default:"localhost"`
	//Port     uint16 `type:"number" default:"5432"`
	//User     string `required:"true"`
	//Password string `required:"true" type:"password"`
	//Database string `required:"true"`
	//Schema   string `default:"public"`
}
type Driver struct{}

//go:embed icon.svg
var logo string

func (d *Driver) ThunderSourceDriver() thunder.SourceDriverInfo {
	return thunder.SourceDriverInfo{
		ID:  "postgresql_trigger",
		New: func() thunder.SourceDriver { return new(Driver) },

		Name:   "PostgreSQL (Trigger)",
		Image:  logo,
		Config: DriverConfig{},
		Notes:  []string{"NOT WORKING (use WAL driver instead)"},
	}
}

func (d *Driver) TestConfig(config any) (string, error) {
	_, ok := config.(*DriverConfig)
	if !ok {
		return "", errors.New("invalid config type")
	}
	return "", errors.New("not implemented")
}

func (d *Driver) Stats(config any) (*thunder.SourceDriverStats, error) {
	return nil, errors.New("not implemented")
}

var (
	_ thunder.SourceDriver = (*Driver)(nil) // Interface implementation
)
