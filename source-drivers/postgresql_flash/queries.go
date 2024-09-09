package postgresql_flash

import "fmt"

func StatsQuery(schema string) string {
	return fmt.Sprintf(`
SELECT
    table_name as name,
    ARRAY_AGG(DISTINCT column_name) AS columns,
    ARRAY_AGG(DISTINCT column_name) FILTER (WHERE primary_key) AS primary_keys
FROM (
         SELECT
             pgc.relname AS table_name,
             a.attname AS column_name,
             COALESCE(i.indisprimary, false) AS primary_key
         FROM
             pg_attribute a
                 JOIN pg_class pgc ON pgc.oid = a.attrelid
                 LEFT JOIN pg_index i ON (pgc.oid = i.indrelid AND a.attnum = ANY(i.indkey))
                 LEFT JOIN pg_catalog.pg_namespace n ON n.oid = pgc.relnamespace
         WHERE
             pgc.relkind IN ('r', '')  -- Relkind for tables
           AND n.nspname = '%s'
           AND a.attnum > 0
           AND pg_table_is_visible(pgc.oid)
           AND NOT a.attisdropped
     ) AS subquery
GROUP BY table_name
ORDER BY table_name
`, schema)
}
