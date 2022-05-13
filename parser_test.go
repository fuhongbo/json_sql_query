package lexerstudy

import (
	"fmt"
	"strings"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	stmt, err := NewParser(strings.NewReader("select name as hello, ba as c,m, 5.8 as x FROM tbl where ba=1 and (c>3 or d.c='ll' or e=3) and v=3")).Parse()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%v", stmt)
	}
}

func BenchmarkNewAST(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewParser(strings.NewReader("select name as hello, ba as c,m FROM tbl where ba=1 and v=3")).Parse()

	}
}
