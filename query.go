package json_sql_query

import (
	"github.com/fuhongbo/json_sql_query/gjson"
	"github.com/tidwall/sjson"
	"strings"
)

type Ref struct {
	Str     string
	JsonRef gjson.Result
}

type Query struct {
	Statement SelectStatement
}

func NewQuery(query string) (*Query, error) {
	selectStatement, err := NewParser(query).Parse()
	if err != nil {
		return nil, err
	} else {
		return &Query{Statement: selectStatement}, nil
	}
}

func (q *Query) Update(query string) error {
	selectStatement, err := NewParser(query).Parse()
	if err != nil {
		return err
	} else {
		q.Statement = selectStatement
	}
	return nil
}

func (q *Query) Valid(json string) (bool, string) {

	ref := &Ref{JsonRef: gjson.Parse(json)}

	result := q.valid(ref)

	if result {
		b := NewBuilder()
		for _, item := range q.Statement.Fields {
			switch item.Type {
			case ALLField:
				return true, ref.JsonRef.Raw
			case IDINTField:
				if item.Alias != "" {
					b.Set(item.Alias, ref.JsonRef.Get(item.Name).Value())
				} else {
					b.Set(item.Name, ref.JsonRef.Get(item.Name).Value())
				}
			case ValueField:
				b.Set(item.Alias, item.Value)
			}
		}
		return result, b.ToJson()

	}
	return result, ""
}

func (q *Query) valid(ref *Ref) bool {

	if q.Statement.WHERE == nil {
		return true
	}
	result, ar := Valid(q.Statement.WHERE.Left, q.Statement.WHERE.Right, q.Statement.WHERE.Op, ref)
	if result && ar != nil {
		key := strings.Split(q.Statement.WHERE.Left.Name, "->")[0]
		tempJson, _ := sjson.Set(ref.JsonRef.Raw, key, ar[key])
		tp := gjson.Parse(tempJson)
		ref.JsonRef = tp
	}
	return result
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
