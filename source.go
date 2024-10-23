package thunder

import (
	"github.com/oklog/ulid/v2"
	"github.com/quix-labs/thunder/utils"
)

type Source struct {
	ID     string
	Driver SourceDriver
	Config utils.DynamicConfig
}

var Sources = utils.NewRegistry[Source]("source").SetIdGenerator(func(item *Source) (string, error) {
	ulid := ulid.Make().String()
	item.ID = ulid
	return ulid, nil
})
