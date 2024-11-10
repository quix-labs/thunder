package elastic

import (
	"fmt"
	"github.com/quix-labs/thunder/utils"
	"strings"
)

type DriverConfig struct {
	Endpoint string `type:"url" required:"true" default:"http://localhost:9200"`
	Username string
	Password string `type:"password"`

	BatchMaxBytesSize int `type:"number" label:"Maximum Batch Size (Kilobytes)" default:"5120" min:"1024" help:"Maximum size of each batch in Kilobytes. (Default to 5MB)"`
	ParallelBatch     int `type:"number" label:"Parallel Batches" default:"4" min:"1" help:"Number of batches to send in parallel. Use 1 to disable parallel processing, which may slow down processing."`

	ReactivityInterval int `type:"number" label:"Reactivity Interval (in sec)" default:"10" min:"1" help:"Max time without data before flushing changes to the index"`
	Prefix             string
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
		cnx += " (prefix:" + cfg.Prefix + ")"
	}
	return cnx
}

var (
	_ utils.DynamicConfig = (*DriverConfig)(nil)
)

func main() {}
