module github.com/quix-labs/thunder/modules/http_server

go 1.23.0

replace github.com/quix-labs/thunder => ../../

replace github.com/quix-labs/thunder/utils => ../../utils

require github.com/quix-labs/thunder v0.0.0-00010101000000-000000000000

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/quix-labs/thunder/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
