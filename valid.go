package json_sql_query

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
	"json_sql_query/gjson"
	"json_sql_query/wildcard"
	"reflect"
	"strings"
)

var floatType = reflect.TypeOf(float64(0))

func Valid(left *Condition, right *Condition, op Token, ref *Ref) (bool, map[string]interface{}) {

	if left == nil && right == nil {
		return true, nil
	} else {
		switch op {
		case AND: // AND
			l, ltr := Valid(left.Left, left.Right, left.Op, ref)
			if l && ltr != nil {
				key := strings.Split(left.Left.Name, "->")[0]
				tempJson, _ := sjson.Set(ref.JsonRef.Raw, key, ltr[key])
				tp := gjson.Parse(tempJson)
				ref.JsonRef = tp
			}
			//在AND中，如果有一个值为false，那么直接返回就行，不需要在向下计算
			if !l {
				return false, nil
			}

			r, rtr := Valid(right.Left, right.Right, right.Op, ref)

			if r && rtr != nil {
				key := strings.Split(right.Left.Name, "->")[0]
				tempJson, _ := sjson.Set(ref.JsonRef.Raw, key, rtr[key])
				tp := gjson.Parse(tempJson)
				ref.JsonRef = tp
			}

			return l && r, nil

		case OR: //  OR
			l, ltr := Valid(left.Left, left.Right, left.Op, ref)
			r, rtr := Valid(right.Left, right.Right, right.Op, ref)

			if ltr != nil && len(ltr) > 0 && rtr != nil && len(rtr) > 0 {

				key := strings.Split(left.Left.Name, "->")[0]
				m := append(ltr[key].([]interface{}), rtr[key].([]interface{}))
				x := removeDuplicateValues(m)
				tempJson, _ := sjson.Set(ref.JsonRef.Raw, key, x)
				tp := gjson.Parse(tempJson)
				ref.JsonRef = tp
			} else {
				if ltr != nil && len(ltr) > 0 {
					key := strings.Split(left.Left.Name, "->")[0]
					tempJson, _ := sjson.Set(ref.JsonRef.Raw, key, ltr[key])
					tp := gjson.Parse(tempJson)
					ref.JsonRef = tp
				}

				if rtr != nil && len(rtr) > 0 {
					key := strings.Split(left.Left.Name, "->")[0]
					tempJson, _ := sjson.Set(ref.JsonRef.Raw, key, rtr[key])
					tp := gjson.Parse(tempJson)
					ref.JsonRef = tp
				}

			}

			return l || r, nil

		default:
			if strings.Contains(left.Name, "->") {
				arrayFilter := make(map[string]interface{})
				return compare(getSpecValue(left, right, op, ref.JsonRef, arrayFilter), nil, op), arrayFilter
			} else {
				return compare(getValue(left, ref.JsonRef), getValue(right, ref.JsonRef), op), nil
			}

		}
	}

}

func compare(a, b interface{}, op Token) bool {

	switch a.(type) {
	case float64:
		tempA := a.(float64)

		switch op {
		case EQ: //   =
			tempB := getFloat(b)
			return tempA == tempB
		case NEQ: //  <>
			tempB := getFloat(b)
			return tempA != tempB
		case GT: //   >
			tempB := getFloat(b)
			return tempA > tempB
		case GTE: //  >=
			tempB := getFloat(b)
			return tempA >= tempB
		case LT: //   <
			tempB := getFloat(b)
			return tempA < tempB
		case LTE: // <=
			tempB := getFloat(b)
			return tempA <= tempB
		case IN:
			return slices.Contains(strings.Split(b.(string), ","), fmt.Sprintf("%v", a))
		case NOTIN:
			return !slices.Contains(strings.Split(b.(string), ","), fmt.Sprintf("%v", a))
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
			return slices.Contains(strings.Split(b.(string), ","), fmt.Sprintf("%v", a))
		case NOTIN:
			return !slices.Contains(strings.Split(b.(string), ","), fmt.Sprintf("%v", a))
		case LIKE:
			return strings.Contains(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
		default:
			return false
		}

	}

}

func getFloat(unk interface{}) float64 {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0
	}
	fv := v.Convert(floatType)
	return fv.Float()
}

func getValue(node *Condition, json gjson.Result) interface{} {
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

func getSpecValue(node *Condition, right *Condition, op Token, json gjson.Result, arrayFilter map[string]interface{}) interface{} {

	c := strings.Replace(node.Name, "->", ".#(", -1)

	rightValue := ""
	switch right.Type {
	case IntegerNode:
		rightValue = fmt.Sprintf("%v", right.IntVal)
	case FloatNode:
		rightValue = fmt.Sprintf("%v", right.FloatVal)
	default:
		rightValue = fmt.Sprintf("%v", right.StrVal)
	}
	switch op {
	case EQ:
		c = c + "==" + rightValue
	case NEQ:
		c = c + "!=" + rightValue
	case LT:
		c = c + "<" + rightValue
	case LTE:
		c = c + "<=" + rightValue
	case GT:
		c = c + ">" + rightValue
	case GTE:
		c = c + ">=" + rightValue
	case IN:
		c = c + "~" + rightValue
	case NOTIN:
		c = c + "~!" + rightValue
	case LIKE:
		c = c + "%" + rightValue
	}
	c = c + ")#"
	tempV := json.Get(c)
	arrayFilter[strings.Split(node.Name, "->")[0]] = tempV.Value()
	return len(tempV.Indexes) > 0
}

func removeDuplicateValues(src []interface{}) []interface{} {

	var result []interface{}
	temp := make(map[string]bool)
	for _, item := range src {
		t, _ := json.Marshal(item)
		if ok := temp[string(t)]; !ok {
			temp[string(t)] = true
			result = append(result, item)
		}
	}
	return result
}
