package json_sql_query

import (
	"fmt"
	"testing"
)

func Test_QueryTravel(t *testing.T) {

	q, err := NewQuery("select a,b,5 as d from tbl where a=1 and b=c")

	if err != nil {
		t.Error(err.Error())
	} else {
		fmt.Println(q.Travel())
	}

}

func Test_QueryData1(t *testing.T) {

	js := `{
    "name": {
        "first": "Tom",
        "last": "Anderson"
    },
    "age": 37,
    "children": [
        "Sara",
        "Alex",
        "Jack"
    ],
    "fav.movie": "Deer Hunter",
    "test": {
        "friends": [
            {
                "first": "Dale",
                "last": "Murphy",
                "age": 44,
                "nets": [
                    "ig",
                    "fb",
                    "tw"
                ]
            },
            {
                "first": "Roger",
                "last": "Craig",
                "age": 68,
                "nets": [
                    "fb",
                    "tw"
                ]
            },
            {
                "first": "Jane",
                "last": "Murphy",
                "age": 47,
                "nets": [
                    "ig",
                    "tw"
                ]
            }
        ]
    }
}`

	q, err := NewQuery(`select name,age,children,test from tbl where test.friends->first not in 'Roger,Jane'`)

	if err != nil {
		t.Error(err.Error())
	} else {
		res, js := q.Valid(js)

		fmt.Println(res)

		fmt.Println(js)
	}

}

func Benchmark_F1(b *testing.B) {
	js := `{
    "name": {
        "first": "Tom",
        "last": "Anderson"
    },
    "age": 37,
    "children": [
        "Sara",
        "Alex",
        "Jack"
    ],
    "fav.movie": "Deer Hunter",
    "test": {
        "friends": [
            {
                "first": "Dale",
                "last": "Murphy",
                "age": 44,
                "nets": [
                    "ig",
                    "fb",
                    "tw"
                ]
            },
            {
                "first": "Roger",
                "last": "Craig",
                "age": 68,
                "nets": [
                    "fb",
                    "tw"
                ]
            },
            {
                "first": "Jane",
                "last": "Murphy",
                "age": 47,
                "nets": [
                    "ig",
                    "tw"
                ]
            }
        ]
    }
}`
	q, err := NewQuery(`select name,age,children,test from tbl where age=37`)

	for i := 0; i < b.N; i++ {

		if err != nil {
			fmt.Println(err.Error())
		} else {
			res, _ := q.Valid(js)
			if !res {
				fmt.Println("粗我")
			}
		}
	}
}
