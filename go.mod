module github.com/quix-labs/thunder

go 1.22

require (
	github.com/creasty/defaults v1.8.0
	github.com/quix-labs/thunder/utils v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace github.com/quix-labs/thunder/utils => ./utils
