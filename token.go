package lexerstudy

type Token int

const (
	ILLEGAL Token = iota
	UNKNOWN_TYPE
	EOF
	WS

	IDENT
	AS

	STRING
	FLOAT
	INTEGER

	ASTERISK //*
	COMMA    //,

	SELECT
	FROM
	WHERE
	AND
	OR

	LT     // <
	LTE    // <=
	NEQ    // <>
	GT     // >
	GTE    // >=
	EQ     // =
	MATCH  //=~
	LEFTC  //(
	RIGTHC //)
)
