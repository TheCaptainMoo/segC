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

	main := FindFunction(functions, "main")

	if main == nil {
		log.Fatalln("ERROR: Entry point not defined!")
	}

	errCode := 0
	if v := main.Body.Interpret(nil); v != nil {
		errCode = int(v.(float64))
	}
	os.Exit(errCode)
}
