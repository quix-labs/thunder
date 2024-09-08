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
	github.com/quix-labs/thunder/source-drivers/postgresql_wal v0.0.0-00010101000000-000000000000
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/doug-martin/goqu/v9 v9.19.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
