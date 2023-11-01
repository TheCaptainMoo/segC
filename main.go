package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalln("ERROR: No input file being passed!")
	}
	fileContent, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatalln(err)
	}

	tokens := Tokenize(string(fileContent))

	ast := ParseProgram(&tokens)

	ast.Interpret(nil)

	main, ok := functions["main"]

	if !ok {
		log.Fatalln("ERROR: Entry point not defined!")
	}

	var exprs []Expression

	if len(args)-2 != len(main.Args) {
		log.Println("ERROR: Number of Command-line arguments not matching!")
		log.Fatalln("Expected: ", len(main.Args))
		log.Fatalln("Recived: ", len(args)-2)
	}

	for i := 2; i < len(args); i++ {
		exprs = append(exprs, ExpressionLiteral{Type: num_literal, Tok: Token{Type: TOK_NUM, Lexme: args[i]}})
	}

	errCode := main.Call(exprs, nil)
	os.Exit(errCode)
}
