package main

import (
	"fmt"
	"log"
	"unicode"
)

const (
	TOK_ID = iota
	TOK_SET
	TOK_FN
	TOK_RET
	TOK_NUM
	TOK_PLUS
	TOK_MINUS
	TOK_MULT
	TOK_DIV
	TOK_EQUALS
	TOK_LESS_THAN
	TOK_GREATER_THAN
	TOK_AND
	TOK_NOT
	TOK_OR
	TOK_IF
	TOK_ELSIF
	TOK_WHILE
	TOK_ELSE
	TOK_OPEN_PARENTH
	TOK_CLOSE_PARENTH
	TOK_OPEN_CURLY
	TOK_CLOSE_CURLY
)

type Token struct {
	Type  uint
	Lexme string
}

func PrintToken(tok Token) {
	fmt.Printf("Token ( %s )\n", tok.Lexme)
}

func TokenTypeFromString(str string) uint {
	switch str {
	case "fn":
		return TOK_FN
	case "ret":
		return TOK_RET
	case "if":
		return TOK_IF
	case "else":
		return TOK_ELSE
	case "elsif":
		return TOK_ELSIF
	case "while":
		return TOK_WHILE
	}
	return TOK_ID
}

func Tokenize(source string) []Token {
	var tokens []Token

	i := 0
	for ; i < len(source); i++ {
		if unicode.IsSpace(rune(source[i])) {
			continue
		}
		switch rune(source[i]) {
		case ';':
			for ; rune(source[i]) != '\n'; i++ {
			}
		case '+':
			tokens = append(tokens, Token{Type: TOK_PLUS, Lexme: "+"})
		case '-':
			tokens = append(tokens, Token{Type: TOK_MINUS, Lexme: "-"})
		case '*':
			tokens = append(tokens, Token{Type: TOK_MULT, Lexme: "*"})
		case '/':
			tokens = append(tokens, Token{Type: TOK_DIV, Lexme: "/"})
		case '=':
			if rune(source[i+1]) == '=' {
				i++
				tokens = append(tokens, Token{Type: TOK_EQUALS, Lexme: "=="})
			} else {
				tokens = append(tokens, Token{Type: TOK_SET, Lexme: "="})
			}
		case '&':
			tokens = append(tokens, Token{Type: TOK_AND, Lexme: "&"})
		case '|':
			tokens = append(tokens, Token{Type: TOK_OR, Lexme: "|"})
		case '<':
			tokens = append(tokens, Token{Type: TOK_LESS_THAN, Lexme: "<"})
		case '>':
			tokens = append(tokens, Token{Type: TOK_GREATER_THAN, Lexme: ">"})
		case '!':
			tokens = append(tokens, Token{Type: TOK_NOT, Lexme: "!"})
		case '(':
			tokens = append(tokens, Token{Type: TOK_OPEN_PARENTH, Lexme: "("})
		case ')':
			tokens = append(tokens, Token{Type: TOK_CLOSE_PARENTH, Lexme: ")"})
		case '{':
			tokens = append(tokens, Token{Type: TOK_OPEN_CURLY, Lexme: "{"})
		case '}':
			tokens = append(tokens, Token{Type: TOK_CLOSE_CURLY, Lexme: "}"})
		default:
			if unicode.IsLetter(rune(source[i])) {
				lexme := ""
				for ; unicode.IsLetter(rune(source[i])) || unicode.IsDigit(rune(source[i])) || rune(source[i]) == '_'; i++ {
					lexme += string(source[i])
				}
				i--
				tokens = append(tokens, Token{Type: TokenTypeFromString(lexme), Lexme: lexme})
			} else if unicode.IsDigit(rune(source[i])) {
				lexme := ""
				for ; unicode.IsDigit(rune(source[i])); i++ {
					lexme += string(source[i])
				}
				i--
				tokens = append(tokens, Token{Type: TOK_NUM, Lexme: lexme})

			} else {
				log.Fatalln("ERROR: Unknown character found! (", string(source[i]), ")")
			}
		}
	}

	return tokens
}
