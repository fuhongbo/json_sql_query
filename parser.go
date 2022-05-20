package json_sql_query

import (
	"fmt"
	"strconv"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r string) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (SelectStatement, error) {
	stmt := SelectStatement{}
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return SelectStatement{}, fmt.Errorf("found %q, expected SELECT", lit)
	}
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK && tok != FLOAT && tok != STRING && tok != INTEGER {
			return SelectStatement{}, fmt.Errorf("found %q, expected field", lit)
		}
		field := Field{}
		switch tok {
		case ASTERISK:
			field.Type = ALLField
			field.Name = lit
		case IDENT:
			field.Type = IDINTField
			field.Name = lit
		case INTEGER:
			field.Type = ValueField
			intValue, _ := strconv.ParseInt(lit, 10, 64)
			field.Value = intValue
		case STRING:
			field.Type = ValueField
			field.Value = lit
		case FLOAT:
			field.Type = ValueField
			floatValue, _ := strconv.ParseFloat(lit, 64)
			field.Value = floatValue
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok == AS {
			_, alias := p.scanIgnoreWhitespace()
			field.Alias = alias

		} else {
			p.unscan()
		}

		stmt.Fields = append(stmt.Fields, field)

		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return SelectStatement{}, fmt.Errorf("found %q, expected FROM", lit)
	}
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT && tok != STRING {
		return SelectStatement{}, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	tok, lit = p.scanIgnoreWhitespace()

	if tok == WHERE {
		ast := NewAST(p)
		if ast.Err != nil {
			fmt.Println("ERROR: " + ast.Err.Error())
		}
		// AST builder
		ar := ast.ParseExpression()
		if ast.Err != nil {
			fmt.Println("ERROR: " + ast.Err.Error())
		}
		stmt.WHERE = ar
	}
	return stmt, nil

}

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}
	tok, lit = p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit
	return
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) unscan() { p.buf.n = 1 }
