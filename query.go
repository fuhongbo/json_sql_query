package json_sql_query

import (
	"json_sql_query/gjson"
	"strings"
)

type Query struct {
	query     string
	Statement *SelectStatement
}

func NewQuery(query string) (*Query, error) {
	selectStatement, err := NewParser(query).Parse()
	if err != nil {
		return nil, err
	} else {
		return &Query{Statement: selectStatement, query: query}, nil
	}
}

func (q *Query) Valid(json string) (bool, string) {

	jsonRef := gjson.Parse(json)

	result, tempArray := q.valid(&jsonRef)

	if result {
		b := NewBuilder()
		for _, item := range q.Statement.Fields {
			switch item.Type {
			case ALLField:
				return true, json
			case IDINTField:
				if item.Alias != "" {
					b.Set(item.Alias, jsonRef.Get(item.Name).Value())
				} else {
					b.Set(item.Name, jsonRef.Get(item.Name).Value())
				}
			case ValueField:
				b.Set(item.Alias, item.Value)
			}
		}

		for k, v := range tempArray {
			b.Set(k, v)
		}

		return result, b.ToJson()

	}
	return result, ""
}

func (e *Query) valid(json *gjson.Result) (bool, map[string]interface{}) {
	if e.Statement.WHERE == nil {
		return true, nil
	}
	tempArray := make(map[string]interface{})
	return Valid(e.Statement.WHERE.Left, e.Statement.WHERE.Right, e.Statement.WHERE.Op, json, tempArray), tempArray
}

func (q *Query) GetFrom() string {
	return q.Statement.TableName
}

func (q *Query) Travel() string {
	travel := "SELECT "
	for _, item := range q.Statement.Fields {
		travel = travel + item.Str() + ","
	}
	travel = strings.TrimRight(travel, ",")
	travel = travel + " FROM " + q.Statement.TableName
	if q.Statement.WHERE != nil {
		travel = travel + " WHERE "
		travel = travel + q.Statement.WHERE.str()
	}

	return travel
}
