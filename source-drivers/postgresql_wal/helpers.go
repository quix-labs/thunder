package postgresql_wal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v5"
	"github.com/quix-labs/thunder"
	"math/rand"
	"strconv"
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
func GetSqlForMapping(table string, mapping *thunder.Mapping) (string, error) {
	// Création de la requête principale
	query := goqu.Dialect("postgres").From(goqu.T(table).As(table))

	var mappingJoins joins
	resultExpr, err := processMapping(table, mapping, &mappingJoins, false)
	if err != nil {
		return "", err
	}

	// Sélection du résultat JSON
	query = query.Select(resultExpr.As("result"))

	// Ajout des jointures dans la requête
	for _, join := range mappingJoins {
		query = query.LeftOuterJoin(join.Table, join.On)
	}

	// Génération de la requête SQL finale
	sql, _, _ := query.ToSQL()
	fmt.Println(sql)
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
		return goqu.Func("json_agg", combined), nil
	}

	return combined, nil
}

//
//func GetSqlForMapping(table string, mapping *thunder.Mapping) (string, error) {
//	query := goqu.Dialect("postgres").From(goqu.T(table).As(table))
//
//	var joins joins
//
//	resultExpr, err := processMapping(table, mapping, &joins, false)
//	if err != nil {
//		return "", err
//	}
//
//	query = query.Select(resultExpr.As("result"))
//
//	for _, join := range joins {
//		query = query.LeftOuterJoin(join.Table, join.On)
//	}
//
//	sql, _, _ := query.ToSQL()
//	fmt.Println(sql)
//	return sql, nil
//}
//
//func processMapping(tableAlias string, mapping *thunder.Mapping, mappingJoins *joins, aggregate bool) (exp.Aliaseable, error) {
//	var jsonSubSelects []any
//	for _, field := range *mapping {
//		if field.FieldType == "simple" {
//			expectedName := field.Column
//			if field.Name != "" {
//				expectedName = field.Name
//			}
//			jsonSubSelects = append(jsonSubSelects, goqu.V(expectedName), goqu.I(tableAlias+"."+field.Column))
//			continue
//		}
//
//		if field.FieldType == "relation" {
//			var relationJoins joins
//
//			relationExpr, err := processMapping(field.Table, &field.Mapping, &relationJoins, field.Type == "has-many")
//			if err != nil {
//				return nil, err
//			}
//
//			// CREATE RELATION QUERY AND ADD NEEDED FIELDS
//			query := goqu.Dialect("postgres").From(goqu.T(field.Table))
//
//			query = query.Select(relationExpr.As("result"))
//			for _, join := range relationJoins {
//				query = query.LeftOuterJoin(join.Table, join.On)
//			}
//
//			relationAlias := field.Table + "_" + generateAlias()
//
//			if !field.UsePivotTable {
//				query = query.SelectAppend(goqu.I(field.Table + "." + field.ForeignKey).As("_pkey"))
//			}
//
//			if !field.UsePivotTable && field.Type == "has-many" {
//				query = query.GroupBy(goqu.I(field.Table + "." + field.ForeignKey))
//			}
//
//			if field.UsePivotTable {
//
//				query = query.From(field.PivotTable).SelectAppend(
//					goqu.I(field.PivotTable+"."+field.ForeignPivotKey).As("_pivot_key"),
//				).InnerJoin(
//					goqu.I(field.Table),
//					goqu.On(goqu.I(field.Table+"."+field.ForeignKey).Eq(goqu.I(field.PivotTable+"."+field.LocalPivotKey))),
//				).GroupBy(goqu.I(field.PivotTable + "." + field.ForeignPivotKey))
//
//				*mappingJoins = append(*mappingJoins, join{
//					query.As(relationAlias),
//					goqu.On(goqu.I(relationAlias + "._pivot_key").Eq(goqu.I(tableAlias + "." + field.LocalKey))),
//				})
//			} else {
//				*mappingJoins = append(*mappingJoins, join{
//					Table: query.As(relationAlias),
//					On:    goqu.On(goqu.I(tableAlias + "." + field.LocalKey).Eq(goqu.I(relationAlias + "._pkey"))),
//				})
//			}
//
//			jsonSubSelects = append(jsonSubSelects, goqu.V(field.Name), goqu.I(relationAlias+".result"))
//		}
//	}
//
//	// Prevent ERROR: cannot pass more than 100 arguments to a function (SQLSTATE 54023) using chunk
//	const chunkSize = 100
//	if len(jsonSubSelects) <= chunkSize {
//		if aggregate {
//			return goqu.Func("json_agg", goqu.Func("json_build_object", jsonSubSelects...)), nil
//		}
//
//		return goqu.Func("json_build_object", jsonSubSelects...), nil
//	}
//
//	var chunks []exp.SQLFunctionExpression
//	for i := 0; i < len(jsonSubSelects); i += chunkSize {
//		end := i + chunkSize
//		if end > len(jsonSubSelects) {
//			end = len(jsonSubSelects)
//		}
//
//		chunks = append(chunks, goqu.Func("jsonb_build_object", jsonSubSelects[i:end]...))
//	}
//
//	// Combine chunks with the || operator
//	var combined exp.Aliaseable
//	if len(chunks) > 0 {
//		combined = chunks[0]
//		for _, chunk := range chunks[1:] {
//			combined = goqu.L(" ? || ? ", combined, chunk)
//		}
//	}
//
//	if aggregate {
//		return goqu.L("", goqu.Func("json_agg", combined)), nil
//	}
//
//	return combined, nil
//}

func GetResultsSync[T any](conn *pgx.Conn, query string, timeout time.Duration) ([]*T, error) {
	var results []*T
	resultsChan, errChan := GetResultsChan[T](conn, query)
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

func GetResultsChan[T any](conn *pgx.Conn, query string) (<-chan *T, <-chan error) {
	query = fmt.Sprintf("SELECT row_to_json(t) FROM (%s) t", query)
	resultsChan := make(chan *T)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultsChan)

		result, err := conn.Query(context.Background(), query)
		if err != nil {
			errorChan <- err
		}
		defer result.Close()

		for result.Next() {
			var rawJson []byte
			if err := result.Scan(&rawJson); err != nil {
				errorChan <- err
			}

			var document T
			if err := json.Unmarshal(rawJson, &document); err != nil {
				errorChan <- err
			}
			resultsChan <- &document
		}
	}()
	return resultsChan, errorChan
}
