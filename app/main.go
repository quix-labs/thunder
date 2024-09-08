package main

import (
	"fmt"
	"github.com/quix-labs/thunder"
	_ "github.com/quix-labs/thunder/modules/api"
	_ "github.com/quix-labs/thunder/modules/frontend"
	"github.com/quix-labs/thunder/modules/http_server"
	_ "github.com/quix-labs/thunder/modules/http_server"
	"time"

	//_ "github.com/quix-labs/thunder/source-drivers/postgresql_trigger"
	_ "github.com/quix-labs/thunder/source-drivers/postgresql_wal"
)

func main() {
	http_server.SetHttpServerAddr(":3002")
	//http_server.SetHttpServerEnabled(false)

	err := thunder.Start()
	if err != nil {
		panic(err)
	}
	thunder.LoadConfig()
	config := thunder.GetConfig()

	processor := config.Processors[2]
	source := config.Sources[processor.Source]
	driver, _ := thunder.GetSourceDriver(source.Driver)

	driverInstance := driver.New()
	//fmt.Println(driverInstance.Stats(source.Config))
	//return
	if err := driverInstance.Start(source.Config); err != nil {
		panic(err)
	}

	resultsChan, errChan := driverInstance.GetDocumentsForProcessor(&processor, 10)
	timeout := time.After(time.Second * 30)
	var docs []*thunder.Document

	for {
		select {
		case <-timeout:
			panic(fmt.Errorf("timeout: no response after 10 seconds"))

		case err := <-errChan:
			panic(err)

		case doc, opened := <-resultsChan:
			if !opened {
				fmt.Printf("%+v\n", docs[0])
				return
			}
			docs = append(docs, doc)
		}
	}
}
