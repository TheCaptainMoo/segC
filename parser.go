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

const (
	num_literal = iota
	id_literal
)

type ExpressionLiteral struct {
	Expression
	Type uint
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

type ExpressionReturn struct {
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

type ExpressionConditional struct {
	Expression
	Tok       Token
	Condition Expression
	Body      ExpressionBlock
	Last      Expression // just for else/elsif statements
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

func TokensShift(tokens *[]Token) Token {
	if len(*tokens) == 0 {
		return Token{}
	}

	tok := (*tokens)[0]
	*tokens = (*tokens)[1:]
	return tok
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

func ParseAssigment(id string, tokens *[]Token) ExpressionAssigment {
	TokensShift(tokens)
	return ExpressionAssigment{Id: id, Value: ParseExpression(tokens)}
}

func ParseBinaryOperation(expr Expression, tokens *[]Token) ExpressionBinOp {
	return ExpressionBinOp{Left: expr, Operation: TokensShift(tokens), Right: ParseExpression(tokens)}
}

func ParseUnaryOperation(tok Token, tokens *[]Token) ExpressionUnaryOp {
	return ExpressionUnaryOp{Value: ParseExpression(tokens), Operation: tok}
}

func ParseConditional(tok Token, tokens *[]Token) ExpressionConditional {
	if tok.Type == TOK_ELSE {
		return ExpressionConditional{Tok: tok, Body: ParseBlock(tokens)}
	}
	return ExpressionConditional{Tok: tok, Condition: ParseExpression(tokens), Body: ParseBlock(tokens)}
}

func ParseCall(id string, tokens *[]Token) ExpressionCall {
	TokensShift(tokens)

	var arguments []Expression

	for (*tokens)[0].Type != TOK_CLOSE_PARENTH {
		arguments = append(arguments, ParseExpression(tokens))
	}
	TokensShift(tokens)
	return ExpressionCall{Id: id, Args: arguments}
}

func ParseExpression(tokens *[]Token) Expression {
	var expr Expression = nil

	tok := TokensShift(tokens)

	switch tok.Type {
	case TOK_RET:
		expr = ParseReturn(tokens)
	case TOK_NUM:
		switch (*tokens)[0].Type {
		case TOK_PLUS, TOK_MINUS, TOK_MULT, TOK_DIV, TOK_EQUALS, TOK_LESS_THAN, TOK_GREATER_THAN, TOK_AND, TOK_OR:
			expr = ParseBinaryOperation(ParseLiteral(tok), tokens)
		default:
			expr = ParseLiteral(tok)
		}
	case TOK_ID:
		switch (*tokens)[0].Type {
		case TOK_SET:
			expr = ParseAssigment(tok.Lexme, tokens)
		case TOK_OPEN_PARENTH:
			expr = ParseCall(tok.Lexme, tokens)
			switch (*tokens)[0].Type {
			case TOK_PLUS, TOK_MINUS, TOK_MULT, TOK_DIV, TOK_EQUALS, TOK_LESS_THAN, TOK_GREATER_THAN, TOK_AND, TOK_OR:
				expr = ParseBinaryOperation(expr, tokens)
			}
		case TOK_PLUS, TOK_MINUS, TOK_MULT, TOK_DIV, TOK_EQUALS, TOK_LESS_THAN, TOK_GREATER_THAN, TOK_AND, TOK_OR:
			expr = ParseBinaryOperation(ParseLiteral(tok), tokens)
		default:
			expr = ParseLiteral(tok)
		}
	case TOK_NOT, TOK_MINUS:
		expr = ParseUnaryOperation(tok, tokens)
	case TOK_IF, TOK_ELSE, TOK_ELSIF, TOK_WHILE:
		expr = ParseConditional(tok, tokens)
	default:
		log.Fatalln("TODO: Implement parsing for token \"" + tok.Lexme + "\" !")
	}
	return expr
}

func ParseReturn(tokens *[]Token) ExpressionReturn {
	return ExpressionReturn{Value: ParseExpression(tokens)}
}

func ParseBlock(tokens *[]Token) ExpressionBlock {
	var block ExpressionBlock

	if v := TokensShift(tokens); v.Type != TOK_OPEN_CURLY {
		fmt.Println("ERROR: Don't forget to open your curly braces when you are supposed to open a code block!\nHere is an example:")
		fmt.Println("Wrong way:\nif x == y\n\tret 1\n}")
		fmt.Println("Right way:\nif x == y {\n\tret 1\n}")
		os.Exit(1)
	}

	indentationLevel := 0
	for (*tokens)[0].Type != TOK_CLOSE_CURLY || indentationLevel != 0 {
		if (*tokens)[0].Type == TOK_OPEN_CURLY {
			indentationLevel++
		} else if (*tokens)[0].Type == TOK_CLOSE_CURLY {
			indentationLevel--
		}

		expr := ParseExpression(tokens)

		if expr != nil {
			if IsType(expr, ExpressionConditional{}) {
				if v := expr.(ExpressionConditional); v.Tok.Type == TOK_ELSE || v.Tok.Type == TOK_ELSIF {
					for i := len(block.Body) - 1; i >= 0; i-- {
						if IsType(block.Body[i], ExpressionConditional{}) {
							v.Last = block.Body[i]
							expr = v
							break
						}
					}

					if expr.(ExpressionConditional).Last == nil {
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
		}

		block.Body = append(block.Body, expr)
	}

	TokensShift(tokens)

	return block
}

func ParseFunctionDeclaration(tokens *[]Token) Statement {
	id := TokensShift(tokens)
	if v := TokensShift(tokens); v.Type != TOK_OPEN_PARENTH {
		fmt.Println("ERROR: Don't forget to open parentheses when declaring a function!\nHere is an example:")
		fmt.Println("Wrong way: fn myFunction x){}")
		fmt.Println("Right way: fn myFunction(x){}")
		os.Exit(1)
	}
	var args []string
	for (*tokens)[0].Type != TOK_CLOSE_PARENTH {
		v := TokensShift(tokens)
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
	TokensShift(tokens)
	body := ParseBlock(tokens)
	return StatementFunctionDeclaration{Id: id.Lexme, Args: args, Body: body}
}

func ParseProgram(tokens *[]Token) StatementProgram {
	var program StatementProgram

	for len(*tokens) > 0 {
		tok := TokensShift(tokens)
		switch tok.Type {
		case TOK_FN:
			program.Body = append(program.Body, ParseFunctionDeclaration(tokens))
		case TOK_ID:
			switch (*tokens)[0].Type {
			case TOK_SET:
				program.Body = append(program.Body, ParseAssigment(tok.Lexme, tokens))
			}
		}
	}
	return program
}