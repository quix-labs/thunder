package mysql

import "fmt"

func StatsQuery(database string) string {
	return fmt.Sprintf(`
SELECT 
    table_name AS name,
    GROUP_CONCAT(DISTINCT column_name ORDER BY column_name) AS columns,
    GROUP_CONCAT(DISTINCT CASE WHEN primary_key THEN column_name END ORDER BY column_name) AS primary_keys
FROM (
    SELECT 
        t.TABLE_NAME AS table_name,
        c.COLUMN_NAME AS column_name,
        CASE WHEN k.CONSTRAINT_NAME = 'PRIMARY' THEN TRUE ELSE FALSE END AS primary_key
    FROM 
        information_schema.COLUMNS c
    JOIN 
        information_schema.TABLES t 
        ON c.TABLE_SCHEMA = t.TABLE_SCHEMA 
        AND c.TABLE_NAME = t.TABLE_NAME
    LEFT JOIN 
        information_schema.KEY_COLUMN_USAGE k 
        ON k.TABLE_SCHEMA = c.TABLE_SCHEMA 
        AND k.TABLE_NAME = c.TABLE_NAME 
        AND k.COLUMN_NAME = c.COLUMN_NAME 
        AND k.CONSTRAINT_NAME = 'PRIMARY'
    WHERE 
        t.TABLE_SCHEMA = '%s'
        AND t.TABLE_TYPE = 'BASE TABLE'
) AS subquery
GROUP BY table_name
ORDER BY table_name;
`, database)
}
