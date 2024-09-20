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
	"reflect"
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
	resultExpr, err := processMapping(processor.Table, &processor.Mapping, &mappingJoins, false, processor.PrimaryKeys)
	if err != nil {
		return "", err
	}

	// Append mapping join, selects
	query = query.Select(resultExpr.As("Json"))
	for _, join := range mappingJoins {
		query = query.LeftOuterJoin(join.Table, join.On)
	}

	// Append conditions
	for _, condition := range processor.Conditions {
		switch condition.Operator {
		case "=":
			query = query.Where(goqu.I(processor.Table + "." + condition.Column).Eq(condition.Value))
			break
		case "is null":
			query = query.Where(goqu.I(processor.Table + "." + condition.Column).IsNull())
			break
		case "is not null":
			query = query.Where(goqu.I(processor.Table + "." + condition.Column).IsNotNull())
			break
		case "is true":
			query = query.Where(goqu.I(processor.Table + "." + condition.Column).IsTrue())
			break
		case "is false":
			query = query.Where(goqu.I(processor.Table + "." + condition.Column).IsFalse())
			break
		default:
			return "", errors.New(fmt.Sprintf("unsupported condition operator: %s", condition.Operator))
		}
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
		goqu.L("?::TEXT", goqu.Func("to_json", goqu.L(fmt.Sprintf("ARRAY[%s]::text[]", strings.Join(bindings, ",")), args...))).As("Pkey"),
	)

	sql, _, _ := query.ToSQL()
	return sql, nil
}

func processMapping(tableAlias string, mapping *thunder.Mapping, mappingJoins *joins, aggregate bool, primaryKeys []string) (exp.Aliaseable, error) {
	var jsonSubSelects []any

	// Append fields
	if mapping.Fields != nil {
		for _, field := range mapping.Fields {
			expectedName := field.Column
			if field.Name != nil && (*field.Name) != "" {
				expectedName = *field.Name
			}
			jsonSubSelects = append(jsonSubSelects, goqu.V(expectedName), goqu.I(tableAlias+"."+field.Column))
		}
	}

	// Recursively append relations
	if mapping.Relations != nil {
		for _, relation := range mapping.Relations {
			var relationJoins joins
			relationExpr, err := processMapping(relation.Table, &relation.Mapping, &relationJoins, relation.Many, relation.PrimaryKeys)
			if err != nil {
				return nil, err
			}

			query := goqu.From(goqu.T(relation.Table)).Select(relationExpr.As("result"))

			for _, join := range relationJoins {
				query = query.LeftOuterJoin(join.Table, join.On)
			}

			relationAlias := relation.Table + "_" + generateAlias()
			if relation.Pivot != nil {
				//TODO PIVOT FIELDS
				pivotAlias := tableAlias + "_" + relation.Pivot.Table + "_" + generateAlias()
				query = query.From(goqu.T(relation.Pivot.Table).As(pivotAlias)).SelectAppend(
					goqu.I(pivotAlias+"."+relation.Pivot.ForeignKey).As("_pivot_key"),
				).InnerJoin(
					goqu.T(relation.Table),
					goqu.On(goqu.I(relation.Table+"."+relation.ForeignKey).Eq(goqu.I(pivotAlias+"."+relation.Pivot.LocalKey))),
				).GroupBy(goqu.I(pivotAlias + "." + relation.Pivot.ForeignKey))

				*mappingJoins = append(*mappingJoins, join{
					Table: query.As(relationAlias),
					On:    goqu.On(goqu.I(relationAlias + "._pivot_key").Eq(goqu.I(tableAlias + "." + relation.LocalKey))),
				})
			} else {
				query = query.SelectAppend(goqu.I(relation.Table + "." + relation.ForeignKey).As("_pkey"))
				if relation.Many {
					query = query.GroupBy(goqu.I(relation.Table + "." + relation.ForeignKey))
				}

				*mappingJoins = append(*mappingJoins, join{
					Table: query.As(relationAlias),
					On:    goqu.On(goqu.I(tableAlias + "." + relation.LocalKey).Eq(goqu.I(relationAlias + "._pkey"))),
				})
			}
			jsonSubSelects = append(jsonSubSelects, goqu.V(relation.Name), goqu.I(relationAlias+".result"))
		}
	}

	// Append primary key columns
	var args []any
	var bindings []string
	for _, primaryColumn := range primaryKeys {
		args = append(args, goqu.I(tableAlias+"."+primaryColumn))
		bindings = append(bindings, "?")
	}
	jsonSubSelects = append(jsonSubSelects,
		goqu.V("_pkey"),
		goqu.L("?::TEXT", goqu.Func("to_json", goqu.L(fmt.Sprintf("ARRAY[%s]::text[]", strings.Join(bindings, ",")), args...))))

	// Convert this to json
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
	resultsChan := make(chan *T)
	errChan := make(chan error, 1)

	go func() {
		defer close(resultsChan)
		GetResultsInChan[T](conn, query, withIntermediateView, resultsChan, errChan)
	}()
	timeoutChan := time.After(timeout)

	for {
		select {
		case <-timeoutChan:
			close(resultsChan)
			close(errChan)
			return nil, fmt.Errorf("timeout: no result received since %s", timeout)

		case err := <-errChan:
			close(errChan)
			close(resultsChan)
			return nil, err

		case row, opened := <-resultsChan:
			if !opened {
				close(errChan)
				return results, nil
			}

			results = append(results, row)
			timeoutChan = time.After(timeout) // RESET TIMEOUT
		}
	}
}
func GetResultsInChan[T any](conn *pgx.Conn, query string, withIntermediateView bool, resultsChan chan<- *T, errorChan chan error) {
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
}

func ExtractKeysFromMapAsJsonString(keys []string, target map[string]any) (string, error) {
	var result = make([]string, len(keys))

	targetType := reflect.TypeOf((*string)(nil)).Elem()

	for idx, key := range keys {
		value, exists := target[key]
		if !exists {
			return "", fmt.Errorf("key %s does not exist", key)
		}

		val := reflect.ValueOf(value)

		if targetType.Kind() == reflect.String {
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				convertedValue := strconv.FormatInt(val.Int(), 10) // Conversion des types entiers signés en string
				result[idx] = any(convertedValue).(string)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				convertedValue := strconv.FormatUint(val.Uint(), 10) // Conversion des types entiers non signés en string
				result[idx] = any(convertedValue).(string)
			case reflect.Float32, reflect.Float64:
				convertedValue := strconv.FormatFloat(val.Float(), 'f', -1, 64) // Conversion des types flottants en string
				result[idx] = any(convertedValue).(string)
			default:
				if val.Type().ConvertibleTo(targetType) {
					convertedValue := val.Convert(targetType).Interface()
					result[idx] = convertedValue.(string)
				} else {
					return "", fmt.Errorf("value for key %s cannot be converted to type %T", key, result[idx])
				}
			}
		} else if val.Type().ConvertibleTo(targetType) {
			convertedValue := val.Convert(targetType).Interface()
			result[idx] = convertedValue.(string)
		} else {
			return "", fmt.Errorf("value for key %s cannot be converted to type %T", key, result[idx])
		}
	}

	return GetPrimaryKeysAsString(result), nil
}

func GetPrimaryKeysAsString(keys []string) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, str := range keys {
		sb.WriteString(strconv.Quote(str))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func MapDiff(map1 map[string]any, map2 map[string]any) map[string]any {
	//return map2
	// TODO PROBABLY USELESS

	diff := make(map[string]any)
	for k, v1 := range map1 {
		if v2, ok := map2[k]; ok {
			if v1 != v2 {
				diff[k] = v2
			}
		}
	}
	for k, v2 := range map2 {
		if _, ok := map1[k]; !ok {
			diff[k] = v2
		}
	}
	return diff
}
