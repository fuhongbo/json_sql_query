package test

import (
	"fmt"
	"github.com/fuhongbo/json_sql_query"
	"testing"
)

const jsonTestStr1 = `{"name":"sqlquery","age":40,"condition":{"height":176,"weight":70,"health":"good"},"friends":[{"name":"Dale","age":44},{"name":"Roger","age":68},{"name":"Jane","age":47}]}`
const jsonTestStr2 = `{"name":"sqlquery","age":40,"topic":"a/b/c","condition":{"height":176,"weight":70,"health":"good"},"friends":[{"name":"Dale","age":44},{"name":"Roger","age":68},{"name":"Jane","age":47}]}`

func TestQuery_Basic_Condition_LT(t *testing.T) {
	sql := "select * from json where age<50"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_LTE(t *testing.T) {
	sql := "select * from json where age<=40"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_NEQ(t *testing.T) {
	sql := "select * from json where age<>40"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, _ := q.Valid(jsonTestStr1)
		if !res {

		} else {
			t.Error("should be can't match")
		}

	}
}

func TestQuery_Basic_Condition_GT(t *testing.T) {
	sql := "select * from json where age>39"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_GTE(t *testing.T) {
	sql := "select * from json where age>=40"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_EQ(t *testing.T) {
	sql := "select * from json where age=40"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_MATCH_PLUS(t *testing.T) {
	sql := "select * from json where topic=~'a/+/c'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr2)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_MATCH_HASHTAG(t *testing.T) {
	sql := "select * from json where topic=~'a/+/c'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr2)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_IN(t *testing.T) {
	sql := "select * from json where age in '40, 41, 42'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_NOTIN(t *testing.T) {
	sql := "select * from json where  age not in '40, 41, 42'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, _ := q.Valid(jsonTestStr1)
		if !res {

		} else {
			t.Error("should be can't match")
		}

	}
}

func TestQuery_Basic_Condition_LIKE(t *testing.T) {
	sql := "select * from json where name like 'query'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_LIKE_NOT_MATCH(t *testing.T) {
	sql := "select * from json where name like 'hello'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, _ := q.Valid(jsonTestStr1)
		if !res {

		} else {
			t.Error("should be can't match")
		}

	}
}

func TestQuery_Basic_Condition_AND(t *testing.T) {
	sql := "select * from json where name='sqlquery' and age=40"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_AND_NOT_MATCH(t *testing.T) {
	sql := "select * from json where name='sqlquery' and age=41"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, _ := q.Valid(jsonTestStr1)
		if !res {

		} else {
			t.Error("should be can't match")
		}

	}
}

func TestQuery_Basic_Condition_OR(t *testing.T) {
	sql := "select * from json where name='sqlquery' or age=41"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_SUB_Filed_Filter(t *testing.T) {

	sql := "select * from json where condition.height>175"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_SUB_Filed_Filter_NOT_MATCH(t *testing.T) {
	sql := "select * from json where condition.height<175"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, _ := q.Valid(jsonTestStr1)
		if !res {

		} else {
			t.Error("should be can't match")
		}

	}
}

func TestQuery_Basic_Condition_Array_Filter(t *testing.T) {

	sql := "select * from json where friends->age >= 47"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_Array_Filter_AND(t *testing.T) {

	sql := "select * from json where friends->age >= 47 and friends->name='Jane'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_Array_Filter_OR(t *testing.T) {

	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health,friends from json where friends->name in 'Roger,Jane' or friends->age<47"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_Array_Filter_IN(t *testing.T) {

	sql := "select * from json where friends->name in 'Roger,Jane'"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_LT_Transform(t *testing.T) {
	sql := "select name,age,condition from json where age<50"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}

func TestQuery_Basic_Condition_LT_Transform_complex(t *testing.T) {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health from json where age<50"
	q, err := json_sql_query.NewQuery(sql)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(jsonTestStr1)
		if !res {
			t.Error("not match")
		}
		fmt.Println(res)

		fmt.Println(js)
	}
}
