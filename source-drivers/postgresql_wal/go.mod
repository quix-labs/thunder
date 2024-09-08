module github.com/quix-labs/thunder/source-drivers/postgresql_wal

go 1.23.0

replace github.com/quix-labs/thunder => ../../

require (
	github.com/jackc/pgconn v1.14.3
	github.com/quix-labs/thunder v0.0.0-00010101000000-000000000000
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/doug-martin/goqu/v9 v9.19.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
