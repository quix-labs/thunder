package postgresql_flash

import (
	"fmt"
	"github.com/quix-labs/thunder/utils"
	"strings"
)

type DriverConfig struct {
	Host     string `default:"localhost"`
	Port     uint16 `type:"number" default:"5432"`
	User     string `required:"true"`
	Password string `required:"true" type:"password"`
	Database string `required:"true"`
	Schema   string `default:"public"`
}

func (cfg DriverConfig) Excerpt() string {
	var cnx string
	cnx = fmt.Sprintf("%s:%s@%s:%d/%s (schema:%s)", cfg.User, strings.Repeat("*", 12), cfg.Host, cfg.Port, cfg.Database, cfg.Schema)
	return cnx
}

var (
	_ utils.DynamicConfig = (*DriverConfig)(nil)
)
