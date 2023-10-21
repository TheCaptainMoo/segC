package main

import (
	"log"
	"strconv"
)

type Variable struct {
	Id    string
	Value float64
}

type Scope struct {
	Variables []Variable
	Parent    *Scope
}

type Function struct {
	Id   string
	Args []string
	Body ExpressionBlock
}

type RetVal struct {
	Value float64
}

var functions []Function
var mainScope Scope = Scope{Parent: nil}

func FindFunction(fns []Function, id string) *Function {
	for i := 0; i < len(fns); i++ {
		if fns[i].Id == id {
			return &fns[i]
		}
	}
	return nil
}

func FindVariable(scope *Scope, id string) *Variable {
	for i := 0; i < len(scope.Variables); i++ {
		if scope.Variables[i].Id == id {
			return &scope.Variables[i]
		}
	}
	if scope.Parent != nil {
		return FindVariable(scope.Parent, id)
	}
	return nil
}

func (n ExpressionBinOp) Interpret(scope *Scope) interface{} {
	if n.Left.Interpret(scope) == nil {
		log.Fatalln("ERROR: Can't operate with a nil type!")
	}
	if n.Right.Interpret(scope) == nil {
		log.Fatalln("ERROR: Can't operate with a nil type!")
	}
	switch n.Operation.Type {
	case TOK_PLUS:
		return n.Left.Interpret(scope).(float64) + n.Right.Interpret(scope).(float64)
	case TOK_MINUS:
		return n.Left.Interpret(scope).(float64) - n.Right.Interpret(scope).(float64)
	case TOK_MULT:
		return n.Left.Interpret(scope).(float64) * n.Right.Interpret(scope).(float64)
	case TOK_DIV:
		return n.Left.Interpret(scope).(float64) / n.Right.Interpret(scope).(float64)
	case TOK_EQUALS:
		if n.Left.Interpret(scope).(float64) == n.Right.Interpret(scope).(float64) {
			return 1.0
		}
		return 0.0
	case TOK_LESS_THAN:
		if n.Left.Interpret(scope).(float64) < n.Right.Interpret(scope).(float64) {
			return 1.0
		}
		return 0.0
	case TOK_GREATER_THAN:
		if n.Left.Interpret(scope).(float64) > n.Right.Interpret(scope).(float64) {
			return 1.0
		}
		return 0.0
	case TOK_AND:
		if n.Left.Interpret(scope).(float64) != 0 && n.Right.Interpret(scope).(float64) != 0 {
			return 1.0
		}
		return 0.0
	case TOK_OR:
		if n.Left.Interpret(scope).(float64) != 0 || n.Right.Interpret(scope).(float64) != 0 {
			return 1.0
		}
		return 0.0
	}
	return nil
}

func (n ExpressionUnaryOp) Interpret(scope *Scope) interface{} {
	if n.Value.Interpret(scope) == nil {
		log.Fatalln("ERROR: Can't operate with a nil type!")
	}
	switch n.Operation.Type {
	case TOK_NOT:
		if n.Value.Interpret(scope).(float64) == 0 {
			return 1.0
		}
		return 0.0
	case TOK_MINUS:
		return -(n.Value.Interpret(scope).(float64))
	}
	return nil
}

func (n ExpressionConditional) Interpret(scope *Scope) interface{} {
	// TODO: Add a check if some idiot puts a "ret" inside a conditional cuz it is 100% posible to do
	switch n.Tok.Type {
	case TOK_IF:
		v := n.Condition.Interpret(scope)
		if v.(float64) != 0 {
			return n.Body.Interpret(scope)
		}
	case TOK_WHILE:
		for n.Condition.Interpret(scope).(float64) != 0 {
			val := n.Body.Interpret(scope)
			if IsType(val, RetVal{}) {
				return val
			}
		}
	case TOK_ELSE:
		v := n.Last.(ExpressionConditional)
		if v.Condition.Interpret(scope).(float64) == 0 {
			if v.Last != nil && v.Last.(ExpressionConditional).Condition.Interpret(scope).(float64) != 0 {
				return nil
			}
			return n.Body.Interpret(scope)
		}
	case TOK_ELSIF:
		v := n.Last.(ExpressionConditional)
		if v.Condition.Interpret(scope).(float64) == 0 && n.Condition.Interpret(scope).(float64) != 0 {
			return n.Body.Interpret(scope)
		}
	}
	return nil
}

func (n StatementProgram) Interpret(scope *Scope) interface{} {
	scope = &mainScope
	for i := 0; i < len(n.Body); i++ {
		n.Body[i].Interpret(scope)
	}
	return nil
}

func (n ExpressionBlock) Interpret(scope *Scope) interface{} {
	if scope == nil {
		scope = &Scope{Parent: &mainScope}
	} else {
		scope = &Scope{Parent: scope}
	}
	for i := 0; i < len(n.Body); i++ {
		result := n.Body[i].Interpret(scope)
		if result != nil {
			if IsType(result, RetVal{}) {
				return result.(RetVal).Value
			}

			if IsType(n.Body[i], ExpressionConditional{}) {
				return result.(float64)
			}
		}
	}
	return nil
}

func (n StatementFunctionDeclaration) Interpret(scope *Scope) interface{} {
	functions = append(functions, Function{Id: n.Id, Args: n.Args, Body: n.Body})
	return nil
}

func (n ExpressionCall) Interpret(scope *Scope) interface{} {
	v := FindFunction(functions, n.Id)
	if v == nil {
		log.Fatalln("ERROR: Can't call function", n.Id, "cuz it is not defined!")
	}
	if len(n.Args) != len(v.Args) {
		log.Fatalln("ERROR: Number of arguments is not matching in call to function", n.Id, "!\nExpected arguments count:", len(v.Args), "\nArguments count recived: ", len(n.Args))
	}

	newScope := Scope{Parent: &mainScope}

	for i := 0; i < len(n.Args); i++ {
		newScope.Variables = append(newScope.Variables, Variable{Id: v.Args[i], Value: n.Args[i].Interpret(scope).(float64)})
	}

	return v.Body.Interpret(&newScope)
}

func (n ExpressionLiteral) Interpret(scope *Scope) interface{} {
	switch n.Type {
	case num_literal:
		val, err := strconv.ParseFloat(n.Tok.Lexme, 64)
		if err != nil {
			log.Fatalln(err)
		}
		return val
	case id_literal:
		v := FindVariable(scope, n.Tok.Lexme)
		if v == nil {
			log.Fatalln("ERROR: Can't access variable", n.Tok.Lexme, "cuz it does not exist!")
		}
		return v.Value
	}
	return nil
}

func (n ExpressionReturn) Interpret(scope *Scope) interface{} {
	return RetVal{Value: n.Value.Interpret(scope).(float64)}
}

func (n ExpressionAssigment) Interpret(scope *Scope) interface{} {
	variable := FindVariable(&mainScope, n.Id)
	if variable == nil {
		mainScope.Variables = append(mainScope.Variables, Variable{Id: n.Id})
		variable = &mainScope.Variables[len(mainScope.Variables)-1]
	}
	v := n.Value.Interpret(scope)
	if v == nil {
		log.Fatalln("ERROR: Can't set variable", n.Id, "to a nil value!")
	}
	variable.Value = v.(float64)
	return variable.Value
}
