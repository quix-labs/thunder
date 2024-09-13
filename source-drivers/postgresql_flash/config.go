package postgresql_flash

import "github.com/quix-labs/thunder"

type DriverConfig struct {
	Host     string `default:"localhost"`
	Port     uint16 `type:"number" default:"5432"`
	User     string `required:"true"`
	Password string `required:"true" type:"password"`
	Database string `required:"true"`
	Schema   string `default:"public"`
}

func (cfg DriverConfig) Excerpt() string {
	return ""
}

var (
	_ thunder.DynamicConfig = (*DriverConfig)(nil)
)
