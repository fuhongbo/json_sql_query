package test

import (
	"github.com/fuhongbo/json_sql_query"
	"testing"
)

func Benchmark_Basic(b *testing.B) {
	sql := "select * from json where age>39"

	for i := 0; i < b.N; i++ {

		q, err := json_sql_query.NewQuery(sql)

		if err != nil {
			b.Error(err.Error())
		} else {
			res, _ := q.Valid(jsonTestStr1)
			if !res {
				b.Error("not match")
			}

		}
	}
}

func Benchmark_Basic_AND_OR(b *testing.B) {
	sql := "select * from json where age>39 and name='sqlquery'"

	for i := 0; i < b.N; i++ {

		q, err := json_sql_query.NewQuery(sql)

		if err != nil {
			b.Error(err.Error())
		} else {
			res, _ := q.Valid(jsonTestStr1)
			if !res {
				b.Error("not match")
			}

		}
	}
}

func Benchmark_Basic_Transform(b *testing.B) {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health from json where age>39 and name='sqlquery'"

	for i := 0; i < b.N; i++ {

		q, err := json_sql_query.NewQuery(sql)

		if err != nil {
			b.Error(err.Error())
		} else {
			res, _ := q.Valid(jsonTestStr1)
			if !res {
				b.Error("not match")
			}

		}
	}
}

func Benchmark_Basic_Array(b *testing.B) {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health from json where friends->name in 'Roger,Jane' "

	for i := 0; i < b.N; i++ {

		q, err := json_sql_query.NewQuery(sql)

		if err != nil {
			b.Error(err.Error())
		} else {
			res, _ := q.Valid(jsonTestStr1)
			if !res {
				b.Error("not match")
			}

		}
	}
}

func Benchmark_Basic_Array_AND(b *testing.B) {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health,friends from json where friends->name in 'Roger,Jane' and friends->age>47 "

	for i := 0; i < b.N; i++ {

		q, err := json_sql_query.NewQuery(sql)

		if err != nil {
			b.Error(err.Error())
		} else {
			res, _ := q.Valid(jsonTestStr1)
			if !res {
				b.Error("not match")
			}

		}
	}
}

func Benchmark_Basic_Array_OR(b *testing.B) {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health from json where friends->name in 'Roger,Jane' or friends->age<47 "

	for i := 0; i < b.N; i++ {

		q, err := json_sql_query.NewQuery(sql)

		if err != nil {
			b.Error(err.Error())
		} else {
			res, _ := q.Valid(jsonTestStr1)
			if !res {
				b.Error("not match")
			}

		}
	}
}
