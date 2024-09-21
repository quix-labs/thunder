module github.com/quix-labs/thunder/modules/frontend

go 1.23.0

replace github.com/quix-labs/thunder => ../../
replace github.com/quix-labs/thunder/utils => ../../utils
replace github.com/quix-labs/thunder/modules/http_server => ../http_server

require (
	github.com/quix-labs/thunder v0.0.0-00010101000000-000000000000
	github.com/quix-labs/thunder/modules/http_server v0.0.0-00010101000000-000000000000
)

require github.com/creasty/defaults v1.8.0 // indirect
