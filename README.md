
<p align="center">
<b style="font-size: 30px">JSON SQL QUERY</b>
<p>
JSON SQL QUERY is a Go package that provides a fast and simple way to filter or transform a json document. You can filter or transform json using SQL-like statements.

Getting Started
===============

## Installing

To start using JSON SQL QUERY, install Go and run `go get`:

```sh
$ go get -u github.com/fuhongbo/json_sql_query
```

This will retrieve the library.

## Filter and transform a json document

Use SQL-like statements to filter transform JSON.

```go
package main

import (
	"fmt"
	"github.com/fuhongbo/json_sql_query"
)

func main() {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health from json where age<50"
	q, err := json_sql_query.NewQuery(sql)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		res, resJson := q.Valid(`{"name":"sqlquery","age":40,"condition":{"height":176,"weight":70,"health":"good"},"friends":[{"name":"Dale","age":44},{"name":"Roger","age":68},{"name":"Jane","age":47}]}`)
		if !res {
			fmt.Println("not match")
		}else{
			fmt.Println(resJson)
		}
	}
}
```

This will print:

```json
{"age":40,"health":"good","height":176,"name":"sqlquery","weight":70}
```

How to filter array in json:
```go
package main

import (
	"fmt"
	"github.com/fuhongbo/json_sql_query"
)

func main() {
	sql := "select name,age,condition.height as height,condition.weight as weight,condition.health as health,friends from json where friends->name in 'Roger,Jane' or friends->age<47"
	q, err := json_sql_query.NewQuery(sql)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		res, resJson := q.Valid(`{"name":"sqlquery","age":40,"condition":{"height":176,"weight":70,"health":"good"},"friends":[{"name":"Dale","age":44},{"name":"Roger","age":68},{"name":"Jane","age":47}]}`)
		if !res {
			fmt.Println("not match")
		}else{
			fmt.Println(resJson)
		}
	}
}
```

This will print:

```json
{"age":40,"friends":[{"age":68,"name":"Roger"},{"age":47,"name":"Jane"},[{"age":44,"name":"Dale"}]],"health":"good","height":176,"name":"sqlquery","weight":70}
```

## Currently, supported conditional expressions 

| Expressions | Remark                               |
|-------------|--------------------------------------|
| =           |                                      |
| <>          |                                      |
| \>          |                                      |
| \>=         |                                      |
| <           |                                      |
| <=          |                                      |
| AND         |                                      |
| OR          |                                      |
| =~          | This special use to MQTT topic match |
| in          |                                      |  
| not in      |                                      |  
| like        |                                      |  
| ()          |                                      |  

## Benchmark

```azure
Benchmark_Basic-16              	  574432	      2105 ns/op	    1528 B/op	      44 allocs/op
Benchmark_Basic_AND_OR-16       	  336627	      3416 ns/op	    2354 B/op	      71 allocs/op
Benchmark_Basic_Transform-16    	  115022	     10074 ns/op	    5826 B/op	     155 allocs/op
Benchmark_Basic_Array-16        	   82503	     15167 ns/op	    8319 B/op	     186 allocs/op
Benchmark_Basic_Array_AND-16    	   49906	     21997 ns/op	   12190 B/op	     261 allocs/op
Benchmark_Basic_Array_OR-16     	   53088	     22883 ns/op	   12326 B/op	     271 allocs/op
```