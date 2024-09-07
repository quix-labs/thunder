module github.com/quix-labs/thunder/app

go 1.23.0

replace github.com/quix-labs/thunder => ../

replace github.com/quix-labs/thunder/source-drivers/postgresql_trigger => ../source-drivers/postgresql_trigger

replace github.com/quix-labs/thunder/source-drivers/postgresql_wal => ../source-drivers/postgresql_wal

replace github.com/quix-labs/thunder/modules/frontend => ../modules/frontend

replace github.com/quix-labs/thunder/modules/http_server => ../modules/http_server

replace github.com/quix-labs/thunder/modules/api => ../modules/api

require (
	github.com/quix-labs/thunder v0.0.0-00010101000000-000000000000
	github.com/quix-labs/thunder/modules/api v0.0.0-00010101000000-000000000000
	github.com/quix-labs/thunder/modules/frontend v0.0.0-00010101000000-000000000000
	github.com/quix-labs/thunder/modules/http_server v0.0.0-00010101000000-000000000000
	github.com/quix-labs/thunder/source-drivers/postgresql_trigger v0.0.0-00010101000000-000000000000
	github.com/quix-labs/thunder/source-drivers/postgresql_wal v0.0.0-00010101000000-000000000000
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
