package json_sql_query

import (
	"fmt"
	"json_sql_query/gjson"
	"json_sql_query/wildcard"
	"reflect"
	"strings"
)

var floatType = reflect.TypeOf(float64(0))

func Valid(left *Condition, right *Condition, op Token, json *gjson.Result, arrayFilter map[string]interface{}) bool {

	if left == nil && right == nil {
		return true
	} else {
		switch op {
		case AND: // AND
			return Valid(left.Left, left.Right, left.Op, json, arrayFilter) && Valid(right.Left, right.Right, right.Op, json, arrayFilter)
		case OR: //  OR
			return Valid(left.Left, left.Right, left.Op, json, arrayFilter) || Valid(right.Left, right.Right, right.Op, json, arrayFilter)
		default:
			if strings.Contains(left.Name, "->") {
				return compare(getSpecValue(left, right, op, json, arrayFilter), getValue(right, json), op)
			} else {
				return compare(getValue(left, json), getValue(right, json), op)
			}

		}
	}

}

func compare(a, b interface{}, op Token) bool {

	switch a.(type) {
	case float64:
		tempA, err := getFloat(a)
		if err != nil {
			return false
		}
		tempB, err := getFloat(b)
		if err != nil {
			return false
		}

		switch op {
		case EQ: //   =
			return tempA == tempB
		case NEQ: //  <>
			return tempA != tempB
		case GT: //   >
			return tempA > tempB
		case GTE: //  >=
			return tempA >= tempB
		case LT: //   <
			return tempA < tempB
		case LTE: // <=
			return tempA <= tempB
		default:
			return false
		}
	case bool:
		return a.(bool)
	default:
		switch op {
		case EQ: //   =
			return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
		case NEQ: //  <>
			return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b)
		case GT: //   >
			return fmt.Sprintf("%v", a) > fmt.Sprintf("%v", b)
		case GTE: //  >=
			return fmt.Sprintf("%v", a) >= fmt.Sprintf("%v", b)
		case LT: //   <
			return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
		case LTE: // <=
			return fmt.Sprintf("%v", a) <= fmt.Sprintf("%v", b)
		case MATCH: // =~
			rs := wildcard.Match(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
			if rs == nil {
				return false
			} else {
				return true
			}
		case IN:
			return true
		case NOTIN:
			return true
		case LIKE:
			return true
		default:
			return false
		}

	}

}

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func getValue(node *Condition, json *gjson.Result) interface{} {
	switch node.Type {
	case FieldNode:
		return json.Get(node.Name).Value()
	case FloatNode:
		return node.FloatVal
	case IntegerNode:
		return node.IntVal
	case StringNode:
		return node.StrVal
	}
	return nil
}

func getSpecValue(node *Condition, right *Condition, op Token, json *gjson.Result, arrayFilter map[string]interface{}) interface{} {

	c := strings.Replace(node.Name, "->", ".#(", -1)

	switch op {
	case EQ:
		c = c + "==" + right.StrVal
	case NEQ:
		c = c + "!=" + right.StrVal
	case LT:
		c = c + "<" + right.StrVal
	case LTE:
		c = c + "<=" + right.StrVal
	case GT:
		c = c + ">" + right.StrVal
	case GTE:
		c = c + ">=" + right.StrVal
	case IN:
		c = c + "~" + right.StrVal
	case NOTIN:
		c = c + "~!" + right.StrVal
	case LIKE:
		c = c + "%" + right.StrVal
	}
	c = c + ")#"
	tempV := json.Get(c)
	arrayFilter[strings.Split(node.Name, "->")[0]] = tempV.Value()
	return len(json.Get(c).Array()) > 0
}
