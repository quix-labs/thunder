package main

import (
	"github.com/quix-labs/thunder"
	_ "github.com/quix-labs/thunder/modules/api"
	_ "github.com/quix-labs/thunder/modules/frontend"
	"github.com/quix-labs/thunder/modules/http_server"
	_ "github.com/quix-labs/thunder/modules/http_server"
	_ "github.com/quix-labs/thunder/source-drivers/postgresql_trigger"
	_ "github.com/quix-labs/thunder/source-drivers/postgresql_wal"
)

func main() {
	http_server.SetHttpServerAddr(":3002")
	//http_server.SetHttpServerEnabled(false)

	err := thunder.Start()
	if err != nil {
		panic(err)
	}
}
