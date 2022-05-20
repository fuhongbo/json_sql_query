package json_sql_query

import (
	"fmt"
	"strconv"
)

type ExprNodeType int
type FieldType int

const (
	Unknown     ExprNodeType = 0
	BinaryNode  ExprNodeType = 1
	FieldNode   ExprNodeType = 2
	IntegerNode ExprNodeType = 3
	FloatNode   ExprNodeType = 4
	StringNode  ExprNodeType = 5
)

const (
	IDINTField FieldType = 0
	ValueField FieldType = 1
	ALLField   FieldType = 2
)

type SelectStatement struct {
	Fields    []Field
	TableName string
	WHERE     *Condition
}

type Field struct {
	Type  FieldType
	Name  string
	Value interface{}
	Alias string
}

type Condition struct {
	Type     ExprNodeType
	Left     *Condition
	Right    *Condition
	FloatVal float64
	IntVal   int64
	StrVal   string
	Name     string
	Op       Token
	OpStr    string
}

func (node *Condition) str() string {
	switch node.Type {
	case FieldNode:
		return node.Name
	case StringNode:
		return node.StrVal
	case IntegerNode:
		return strconv.FormatInt(node.IntVal, 10)
	case FloatNode:
		return strconv.FormatFloat(node.FloatVal, 'f', -1, 64)
	case BinaryNode:
		return fmt.Sprintf("(%s %s %s)", node.Left.str(), node.OpStr, node.Right.str())
	}

	return "?"
}

func (field Field) Str() string {
	f := field.Name
	if field.Type == ValueField {
		f = fmt.Sprintf("%v", field.Value)
	}

	if field.Type == ALLField {
		f = "*"
	}

	if field.Alias != "" {
		f = f + " AS " + field.Alias
	}
	return f
}
