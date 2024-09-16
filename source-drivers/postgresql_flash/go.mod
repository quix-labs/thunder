module github.com/quix-labs/thunder/source-drivers/postgresql_flash

go 1.23.0

replace github.com/quix-labs/thunder => ../../

require (
	github.com/doug-martin/goqu/v9 v9.19.0
	github.com/jackc/pgx/v5 v5.7.0
	github.com/quix-labs/flash v0.0.0-20240715105610-d3fe20e10139
	github.com/quix-labs/flash/drivers/trigger v0.0.0-20240715105610-d3fe20e10139
	github.com/quix-labs/thunder v0.0.0-00010101000000-000000000000
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
