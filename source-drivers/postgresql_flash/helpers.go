package postgresql_flash

import (
	"context"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v5"
	"github.com/quix-labs/thunder"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type join struct {
	Table exp.Expression
	On    exp.JoinCondition
}
type joins []join

func generateAlias() string {
	rand.Seed(time.Now().UnixNano())
	return "t" + strconv.Itoa(rand.Intn(10000)) // Exemple : t1234
}
func GetSqlForProcessor(processor *thunder.Processor) (string, error) {
	query := goqu.Dialect("postgres").From(goqu.T(processor.Table))

	var mappingJoins joins
	resultExpr, err := processMapping(processor.Table, &processor.Mapping, &mappingJoins, false)
	if err != nil {
		return "", err
	}

	// Append mapping join, selects
	query = query.Select(resultExpr.As("Json"))
	for _, join := range mappingJoins {
		query = query.LeftOuterJoin(join.Table, join.On)
	}

	// Append primary keys
	if len(processor.PrimaryKeys) == 0 {
		return "", errors.New("primary keys must be set to fetch documents")
	}
	var args []any
	var bindings []string
	for _, primaryColumn := range processor.PrimaryKeys {
		args = append(args, goqu.I(processor.Table+"."+primaryColumn))
		bindings = append(bindings, "?")
	}
	query = query.SelectAppend(
		goqu.L(fmt.Sprintf("ARRAY[%s]::text[]", strings.Join(bindings, ",")), args...).As("PrimaryKeys"),
	)

	sql, _, _ := query.ToSQL()
	return sql, nil
}

func processMapping(tableAlias string, mapping *thunder.Mapping, mappingJoins *joins, aggregate bool) (exp.Aliaseable, error) {
	var jsonSubSelects []any

	for _, field := range *mapping {
		if field.FieldType == "simple" {
			expectedName := field.Name
			if expectedName == "" {
				expectedName = field.Column
			}
			jsonSubSelects = append(jsonSubSelects, goqu.V(expectedName), goqu.I(tableAlias+"."+field.Column))
			continue
		}

		if field.FieldType == "relation" {
			var relationJoins joins
			relationExpr, err := processMapping(field.Table, &field.Mapping, &relationJoins, field.Type == "has-many")
			if err != nil {
				return nil, err
			}

			query := goqu.From(goqu.T(field.Table)).Select(relationExpr.As("result"))

			for _, join := range relationJoins {
				query = query.LeftOuterJoin(join.Table, join.On)
			}

			relationAlias := field.Table + "_" + generateAlias()
			if field.UsePivotTable {
				pivotAlias := tableAlias + "_" + field.PivotTable + "_" + generateAlias()
				query = query.From(goqu.T(field.PivotTable).As(pivotAlias)).SelectAppend(
					goqu.I(pivotAlias+"."+field.ForeignPivotKey).As("_pivot_key"),
				).InnerJoin(
					goqu.T(field.Table),
					goqu.On(goqu.I(field.Table+"."+field.ForeignKey).Eq(goqu.I(pivotAlias+"."+field.LocalPivotKey))),
				).GroupBy(goqu.I(pivotAlias + "." + field.ForeignPivotKey))

				*mappingJoins = append(*mappingJoins, join{
					Table: query.As(relationAlias),
					On:    goqu.On(goqu.I(relationAlias + "._pivot_key").Eq(goqu.I(tableAlias + "." + field.LocalKey))),
				})
			} else {
				query = query.SelectAppend(goqu.I(field.Table + "." + field.ForeignKey).As("_pkey"))
				if field.Type == "has-many" {
					query = query.GroupBy(goqu.I(field.Table + "." + field.ForeignKey))
				}

				*mappingJoins = append(*mappingJoins, join{
					Table: query.As(relationAlias),
					On:    goqu.On(goqu.I(tableAlias + "." + field.LocalKey).Eq(goqu.I(relationAlias + "._pkey"))),
				})
			}
			jsonSubSelects = append(jsonSubSelects, goqu.V(field.Name), goqu.I(relationAlias+".result"))
		}
	}

	if len(jsonSubSelects) <= 100 {
		if aggregate {
			return goqu.Func("json_agg", goqu.Func("json_build_object", jsonSubSelects...)), nil
		}
		return goqu.Func("json_build_object", jsonSubSelects...), nil
	}

	// Prevent ERROR: cannot pass more than 100 arguments to a function (SQLSTATE 54023) using chunk
	const chunkSize = 100
	var chunks []exp.SQLFunctionExpression

	for i := 0; i < len(jsonSubSelects); i += chunkSize {
		end := i + chunkSize
		if end > len(jsonSubSelects) {
			end = len(jsonSubSelects)
		}
		chunks = append(chunks, goqu.Func("jsonb_build_object", jsonSubSelects[i:end]...))
	}

	// Combine chunks with the || operator
	var combined exp.Aliaseable
	if len(chunks) > 0 {
		combined = chunks[0]
		for _, chunk := range chunks[1:] {
			combined = goqu.L(" ? || ? ", combined, chunk)
		}
	}

	if aggregate {
		return goqu.Func("jsonb_agg", combined), nil
	}

	return combined, nil
}

func GetResultsSync[T any](conn *pgx.Conn, query string, timeout time.Duration, withIntermediateView bool) ([]*T, error) {
	var results []*T
	resultsChan, errChan := GetResultsChan[T](conn, query, withIntermediateView)
	timeoutChan := time.After(timeout)
	for {
		select {
		case <-timeoutChan:
			return nil, fmt.Errorf("timeout: no result received since %s", timeout)

		case err := <-errChan:
			return nil, err

		case row, opened := <-resultsChan:
			if !opened {
				return results, nil
			}

			results = append(results, row)
			timeoutChan = time.After(timeout) // RESET TIMEOUT
		}
	}
}

func GetResultsChan[T any](conn *pgx.Conn, query string, withIntermediateView bool) (<-chan *T, <-chan error) {
	resultsChan := make(chan *T)
	errorChan := make(chan error, 1)
	go func() {
		if withIntermediateView {
			// CAN BE USEFUL IF USER WANT A PREVIEW
			// PREVIOUS TEST PROVE THAT IS SLOWER THAN QUERYING DIRECTLY
			viewName := "tmp_" + generateAlias()
			materializedQuery := fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s", viewName, query)

			_, err := conn.Exec(context.Background(), materializedQuery)
			if err != nil {
				errorChan <- err
			}

			defer func() {
				_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP MATERIALIZED VIEW IF EXISTS %s", viewName))
				if err != nil {
					errorChan <- err
				}
			}()
			query = "SELECT * FROM " + viewName
		}
		result, err := conn.Query(context.Background(), query)
		if err != nil {
			errorChan <- err
		}
		defer result.Close()

		for result.Next() {
			document, err := pgx.RowToStructByName[T](result)
			if err != nil {
				errorChan <- err
			}
			resultsChan <- &document
		}
		close(resultsChan)
	}()
	return resultsChan, errorChan
}
