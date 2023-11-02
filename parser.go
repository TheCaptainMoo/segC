package main

import (
	"fmt"
	"log"
	"os"
)

type ASTNode interface {
	Interpret(scope *Scope) interface{}
}

type Statement interface {
	ASTNode
}

type Expression interface {
	ASTNode
}

type LiteralType int

const (
	num_literal LiteralType = iota
	id_literal
)

type ExpressionLiteral struct {
	Expression
	Type LiteralType
	Tok  Token
}

type ExpressionAssigment struct {
	Expression
	Id    string
	Value Expression
}

type ExpressionCall struct {
	Expression
	Id   string
	Args []Expression
}

type StatementReturn struct {
	Expression
	Value Expression
}

type ExpressionBlock struct {
	Expression
	Body []ASTNode
}

type ExpressionBinOp struct {
	Expression
	Left      Expression
	Operation Token
	Right     Expression
}
type ExpressionUnaryOp struct {
	Expression
	Value     Expression
	Operation Token
}

type StatementConditional struct {
	Statement
	Tok       Token
	Condition Expression
	Body      ExpressionBlock
	Next      Statement // just for else/elsif statements
}

type StatementProgram struct {
	Statement
	Body []ASTNode
}

type StatementFunctionDeclaration struct {
	Statement
	Id   string
	Args []string
	Body ExpressionBlock
}

func (arr *TokenList) shift() Token {
	if len(*arr) == 0 {
		return Token{}
	}

	v := (*arr)[0]
	*arr = (*arr)[1:]
	return v
}

func (arr *TokenList) cur() Token {
	if len(*arr) == 0 {
		return Token{}
	}

	return (*arr)[0]
}

func ParseLiteral(token Token) ExpressionLiteral {
	switch token.Type {
	case TOK_NUM:
		return ExpressionLiteral{Type: num_literal, Tok: token}

	case TOK_ID:
		return ExpressionLiteral{Type: id_literal, Tok: token}
	}
	return ExpressionLiteral{}
}

func ParseAssigment(id string, tokens *TokenList) ExpressionAssigment {
	tokens.shift()
	return ExpressionAssigment{Id: id, Value: ParseExpression(tokens)}
}

func ParseBinaryOperation(expr Expression, tokens *TokenList) ExpressionBinOp {
	return ExpressionBinOp{Left: expr, Operation: tokens.shift(), Right: ParseExpression(tokens)}
}

func ParseUnaryOperation(tok Token, tokens *TokenList) ExpressionUnaryOp {
	return ExpressionUnaryOp{Value: ParseExpression(tokens), Operation: tok}
}

func ParseConditional(tok Token, tokens *TokenList) StatementConditional {
	if tok.Type == TOK_ELSE {
		return StatementConditional{Tok: tok, Body: ParseBlock(tokens)}
	}
	expr := StatementConditional{Tok: tok, Condition: ParseExpression(tokens), Body: ParseBlock(tokens)}
	if v := tokens.cur(); v.Type == TOK_ELSE || v.Type == TOK_ELSIF {
		expr.Next = ParseConditional(tokens.shift(), tokens)
	}
	return expr
}

func ParseCall(id string, tokens *TokenList) ExpressionCall {
	tokens.shift()
	var arguments []Expression

	for tokens.cur().Type != TOK_CLOSE_PARENTH {
		arguments = append(arguments, ParseExpression(tokens))
	}
	tokens.shift()
	return ExpressionCall{Id: id, Args: arguments}
}

func Parse(tokens *TokenList) ASTNode {
	// Try and parse everything as a Statement
	switch tokens.cur().Type {
	case TOK_RET:
		tokens.shift()
		return ParseReturn(tokens)
	case TOK_IF, TOK_ELSE, TOK_ELSIF, TOK_WHILE:
		return ParseConditional(tokens.shift(), tokens)
	default:
		// If you can't, parse it as an Expression
		return ParseExpression(tokens)
	}
}

func ParseExpression(tokens *TokenList) Expression {
	var expr Expression = nil

	tok := tokens.shift()

	switch tok.Type {
	case TOK_NUM:
		expr = ParseLiteral(tok)
	case TOK_ID:
		switch tokens.cur().Type {
		case TOK_SET:
			expr = ParseAssigment(tok.Lexme, tokens)
		case TOK_OPEN_PARENTH:
			expr = ParseCall(tok.Lexme, tokens)
			switch tokens.cur().Type {
			case TOK_PLUS, TOK_MINUS, TOK_MULT, TOK_DIV, TOK_EQUALS, TOK_LESS_THAN, TOK_GREATER_THAN, TOK_AND, TOK_OR:
				expr = ParseBinaryOperation(expr, tokens)
			}
		default:
			expr = ParseLiteral(tok)
		}
	case TOK_NOT, TOK_MINUS:
		expr = ParseUnaryOperation(tok, tokens)
	default:
		log.Fatalln("TODO: Implement parsing for token \"" + tok.Lexme + "\" !")
	}
	switch tokens.cur().Type {
	case TOK_PLUS, TOK_MINUS, TOK_MULT, TOK_DIV, TOK_EQUALS, TOK_LESS_THAN, TOK_GREATER_THAN, TOK_AND, TOK_OR:
		expr = ParseBinaryOperation(expr, tokens)
	}
	return expr
}

func ParseReturn(tokens *TokenList) StatementReturn {
	return StatementReturn{Value: ParseExpression(tokens)}
}

func ParseBlock(tokens *TokenList) ExpressionBlock {
	var block ExpressionBlock

	if v := tokens.shift(); v.Type != TOK_OPEN_CURLY {
		fmt.Println("ERROR: Don't forget to open your curly braces when you are supposed to open a code block!\nHere is an example:")
		fmt.Println("Wrong way:\nif x == y\n\tret 1\n}")
		fmt.Println("Right way:\nif x == y {\n\tret 1\n}")
		os.Exit(1)
	}

	indentationLevel := 0
	for tokens.cur().Type != TOK_CLOSE_CURLY || indentationLevel != 0 {
		if tokens.cur().Type == TOK_OPEN_CURLY {
			indentationLevel++
		} else if tokens.cur().Type == TOK_CLOSE_CURLY {
			indentationLevel--
		}

		expr := Parse(tokens)

		if expr != nil {
			if IsType(expr, StatementConditional{}) {
				if v := expr.(StatementConditional); v.Tok.Type != TOK_IF || v.Tok.Type == TOK_WHILE {
					fmt.Println("ERROR: Don't forget that", v.Tok.Lexme, "is used after another condition!\nHere is an example:")
					fmt.Println("Wrong way:")
					fmt.Println("fn main() {")
					if v.Tok.Type != TOK_ELSE {
						fmt.Println("\t" + v.Tok.Lexme + " y == 24 {")
						fmt.Println("\t\tret 0")
						fmt.Println("\t}")
					} else {
						fmt.Println("\t"+v.Tok.Lexme, "{")
						fmt.Println("\t\tret 0")
						fmt.Println("\t}")
					}
					fmt.Println("}")
					fmt.Println("Right way:")
					fmt.Println("fn main() {")
					fmt.Println("\tif y != 0 {")
					fmt.Println("\t\tret 1")
					fmt.Println("\t}")
					if v.Tok.Type != TOK_ELSE {
						fmt.Println("\t" + v.Tok.Lexme + " y == 24 {")
						fmt.Println("\t\tret 0")
						fmt.Println("\t}")
					} else {
						fmt.Println("\t"+v.Tok.Lexme, "{")
						fmt.Println("\t\tret 0")
						fmt.Println("\t}")
					}
					fmt.Println("}")
					os.Exit(1)
				}
			}
		}

		block.Body = append(block.Body, expr)
	}

	tokens.shift()

	return block
}

func ParseFunctionDeclaration(tokens *TokenList) Statement {
	id := tokens.shift()
	if v := tokens.shift(); v.Type != TOK_OPEN_PARENTH {
		fmt.Println("ERROR: Don't forget to open parentheses when declaring a function!\nHere is an example:")
		fmt.Println("Wrong way: fn myFunction x){}")
		fmt.Println("Right way: fn myFunction(x){}")
		os.Exit(1)
	}
	var args []string
	for tokens.cur().Type != TOK_CLOSE_PARENTH {
		v := tokens.shift()
		if v.Type == TOK_OPEN_CURLY {
			fmt.Println("ERROR: Don't forget to close parentheses when declaring a function!\nHere is an example:")
			fmt.Println("Wrong way: fn myFunction(x y z {}")
			fmt.Println("Right way: fn myFunction(x y z){}")
			os.Exit(1)
		}
		if v.Type != TOK_ID {
			fmt.Println("ERROR: in function definition, functions arguments should be just identifiers!")
			fmt.Println("Here is an example:\n\tfn myFunction(bar baz)")
			os.Exit(1)
		}
		args = append(args, v.Lexme)
	}
	tokens.shift()
	body := ParseBlock(tokens)
	return StatementFunctionDeclaration{Id: id.Lexme, Args: args, Body: body}
}

func ParseProgram(tokens *TokenList) StatementProgram {
	var program StatementProgram

	for len(*tokens) > 0 {
		tok := tokens.shift()
		switch tok.Type {
		case TOK_FN:
			program.Body = append(program.Body, ParseFunctionDeclaration(tokens))
		case TOK_ID:
			switch tokens.cur().Type {
			case TOK_SET:
				program.Body = append(program.Body, ParseAssigment(tok.Lexme, tokens))
			}
		}
	}
	return program
}
