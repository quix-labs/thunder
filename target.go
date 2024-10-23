package thunder

import (
	"github.com/oklog/ulid/v2"
	"github.com/quix-labs/thunder/utils"
)

var Targets = utils.NewRegistry[Target]("target").SetIdGenerator(func(item *Target) (string, error) {
	ulid := ulid.Make().String()
	item.ID = ulid
	return ulid, nil
})

type Target struct {
	ID     string
	Driver TargetDriver
	Config utils.DynamicConfig
}
