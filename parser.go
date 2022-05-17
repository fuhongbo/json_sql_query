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

func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK && tok != FLOAT && tok != STRING && tok != INTEGER {
			return nil, fmt.Errorf("found %q, expected field", lit)
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
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
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
