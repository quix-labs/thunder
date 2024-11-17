package mysql

import (
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
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

var dialect = "mysql"

func generateAlias() string {
	rand.Seed(time.Now().UnixNano())
	return "t" + strconv.Itoa(rand.Intn(10000)) // Exemple : t1234
}
func GetSqlForProcessor(processor *thunder.Processor) (string, error) {
	query := goqu.Dialect(dialect).From(goqu.T(processor.Table))

	var mappingJoins joins
	resultExpr, err := processMapping(processor.Table, &processor.Mapping, &mappingJoins, false, processor.PrimaryKeys)
	if err != nil {
		return "", err
	}

	// Append mapping join, selects
	query = query.Select(resultExpr.As("json"))
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
		goqu.Cast(
			goqu.Func(
				"JSON_ARRAY",
				func() []interface{} {
					castedBindings := make([]interface{}, len(args))
					for i, arg := range args {
						castedBindings[i] = goqu.Cast(arg.(goqu.Expression), "CHAR")
					}
					return castedBindings
				}()...,
			),
			"CHAR",
		).As("pkey"),

		goqu.Select(
			goqu.Cast(
				goqu.L("1000 * ?", goqu.Func("UNIX_TIMESTAMP", goqu.Func("current_timestamp", 3))),
				"UNSIGNED INTEGER",
			),
		).As("version"),
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

			query := goqu.Dialect(dialect).From(goqu.T(relation.Table)).Select(relationExpr.As("result"))

			for _, join := range relationJoins {
				query = query.LeftOuterJoin(join.Table, join.On)
			}

			relationAlias := relation.Table + "_" + generateAlias()
			if relation.Pivot != nil {
				//TODO PIVOT FIELDS
				pivotAlias := tableAlias + "_" + relation.Pivot.Table + "_" + generateAlias()
				query = query.From(goqu.T(relation.Pivot.Table).As(pivotAlias)).SelectAppend(
					goqu.I(pivotAlias+"."+relation.Pivot.LocalKey).As("_pivot_key"),
				).InnerJoin(
					goqu.T(relation.Table),
					goqu.On(goqu.I(relation.Table+"."+relation.ForeignKey).Eq(goqu.I(pivotAlias+"."+relation.Pivot.ForeignKey))),
				).GroupBy(goqu.I(pivotAlias + "." + relation.Pivot.LocalKey))

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
	var args []goqu.Expression
	for _, primaryColumn := range primaryKeys {
		args = append(args, goqu.I(tableAlias+"."+primaryColumn))
	}

	jsonSubSelects = append(jsonSubSelects,
		goqu.V("_pkey"),
		goqu.Cast(
			goqu.Func(
				"JSON_ARRAY",
				func() []interface{} {
					castedBindings := make([]interface{}, len(args))
					for i, arg := range args {
						castedBindings[i] = goqu.Cast(arg, "CHAR")
					}
					return castedBindings
				}()...,
			),
			"CHAR",
		),
	)

	// Convert this to json
	if aggregate {
		return goqu.Func("JSON_ARRAYAGG", goqu.Func("JSON_OBJECT", jsonSubSelects...)), nil
	}
	return goqu.Func("JSON_OBJECT", jsonSubSelects...), nil
}
