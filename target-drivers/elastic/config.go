package elastic

import (
	"fmt"
	"github.com/quix-labs/thunder"
	"strings"
)

type DriverConfig struct {
	Endpoint  string `default:"http://localhost:9200"`
	Username  string
	Password  string `type:"password"`
	BatchSize int    `type:"number" label:"Batch size" default:"100" min:"1" help:"Use 1 to disable batching (not recommended)"`
	Prefix    string
}

func (cfg DriverConfig) Excerpt() string {
	var cnx string
	if cfg.Username != "" {
		if cfg.Password != "" {
			cnx = fmt.Sprintf("%s:%s@%s", cfg.Username, strings.Repeat("*", len(cfg.Password)), cfg.Endpoint)
		} else {
			cnx = fmt.Sprintf("%s@%s", cfg.Username, cfg.Endpoint)
		}
	} else {
		cnx = cfg.Endpoint
	}

	if cfg.Prefix != "" {
		cnx += "/" + cfg.Prefix
	}
	return cnx
}

var (
	_ thunder.DynamicConfig = (*DriverConfig)(nil)
)