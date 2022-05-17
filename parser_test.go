package json_sql_query

import (
	"fmt"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	stmt, err := NewParser("select * from tbl where b=2 and c like '4' and d->a.id>0").Parse()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%v", stmt)
	}
}

func BenchmarkNewAST(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewParser("select b,c from tbl where b=2").Parse()

	}
}
