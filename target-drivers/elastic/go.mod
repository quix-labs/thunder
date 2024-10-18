module github.com/quix-labs/thunder/target-drivers/elastic

go 1.23.0

replace github.com/quix-labs/thunder => ../../

replace github.com/quix-labs/thunder/utils => ../../utils

require (
	github.com/elastic/go-elasticsearch/v8 v8.15.0
	github.com/quix-labs/thunder v0.0.0-00010101000000-000000000000
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.6.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/quix-labs/thunder/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
