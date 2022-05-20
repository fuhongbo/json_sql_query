package json_sql_query

import (
	"errors"
	"strconv"
	"strings"
)

var precedence = map[string]int{">": 20, ">=": 20, "<": 20, "<=": 20, "<>": 20, "=": 20, "=~": 20, "IN": 20, "NOT IN": 20, "LIKE": 20, "AND": 10, "OR": 10}

type ASTToken struct {
	Tok string
	Token
}

type AST struct {
	p *Parser

	currTok   *ASTToken
	currIndex int
	depth     int

	Err error
}

func NewAST(p *Parser) *AST {
	a := &AST{
		p: p,
	}

	tok, lit := a.p.scanIgnoreWhitespace()

	if tok == EOF {
		a.Err = errors.New("empty token")
	} else {
		a.currIndex = 0
		a.currTok = &ASTToken{
			Tok:   lit,
			Token: tok,
		}
	}

	return a
}

func (a *AST) ParseExpression() *Condition {
	a.depth++ // called depth
	lhs := a.parsePrimary()
	r := a.parseBinOpRHS(0, lhs)
	a.depth--
	if a.depth == 0 && a.Err != nil {
		a.Err = errors.New("error")
	}
	return r
}

func (a *AST) getNextToken() *ASTToken {
	a.currIndex++
	tok, lit := a.p.scanIgnoreWhitespace()

	if tok == EOF {
		return nil
	} else {
		a.currTok = &ASTToken{
			Tok:   lit,
			Token: tok,
		}
		return a.currTok
	}

}

func (a *AST) parsePrimary() *Condition {
	switch a.currTok.Token {
	case IDENT, ILLEGAL:
		c := &Condition{
			Op:   a.currTok.Token,
			Type: FieldNode,
			Name: a.currTok.Tok,
		}
		a.getNextToken()
		return c
	case INTEGER:
		intValue, _ := strconv.ParseInt(a.currTok.Tok, 10, 32)
		c := &Condition{
			Op:     a.currTok.Token,
			Type:   IntegerNode,
			IntVal: intValue,
		}
		a.getNextToken()
		return c
	case FLOAT:
		floatValue, _ := strconv.ParseFloat(a.currTok.Tok, 32)
		c := &Condition{
			Op:       a.currTok.Token,
			Type:     FloatNode,
			FloatVal: floatValue,
		}
		a.getNextToken()
		return c
	case STRING:
		c := &Condition{
			Op:     a.currTok.Token,
			Type:   StringNode,
			StrVal: a.currTok.Tok,
		}
		a.getNextToken()
		return c
	case LEFTC, RIGTHC:

		t := a.getNextToken()
		if t == nil {
			a.Err = errors.New("error")
			return nil
		}
		e := a.ParseExpression()
		if e == nil {
			return nil
		}
		if a.currTok.Tok != ")" {
			a.Err = errors.New("error")
			return nil
		}
		a.getNextToken()
		return e

	default:
		return nil
	}
}

func (a *AST) parseBinOpRHS(execPrec int, lhs *Condition) *Condition {
	for {

		tokPrec := a.getTokPrecedence()
		if tokPrec < execPrec {
			return lhs
		}

		binOp := a.currTok.Tok
		op := a.currTok.Token
		if a.getNextToken() == nil {
			return lhs
		}
		rhs := a.parsePrimary()
		if rhs == nil {
			return nil
		}
		nextPrec := a.getTokPrecedence()
		if tokPrec < nextPrec {
			rhs = a.parseBinOpRHS(tokPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = &Condition{
			Type:  BinaryNode,
			OpStr: binOp,
			Op:    op,
			Left:  lhs,
			Right: rhs,
		}
	}
}

func (a *AST) getTokPrecedence() int {
	if p, ok := precedence[strings.ToUpper(a.currTok.Tok)]; ok {
		return p
	}
	return -1
}
