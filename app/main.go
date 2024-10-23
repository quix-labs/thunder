package main

import (
	"github.com/quix-labs/thunder"
	_ "github.com/quix-labs/thunder/exporters/csv"
	_ "github.com/quix-labs/thunder/exporters/json"
	_ "github.com/quix-labs/thunder/exporters/yaml"
	_ "github.com/quix-labs/thunder/modules/api"
	_ "github.com/quix-labs/thunder/modules/frontend"
	_ "github.com/quix-labs/thunder/modules/http_server"
	_ "github.com/quix-labs/thunder/source-drivers/postgresql_flash"
	_ "github.com/quix-labs/thunder/target-drivers/elastic"
)

func main() {
	err := thunder.Start()
	if err != nil {
		panic(err)
	}
}
