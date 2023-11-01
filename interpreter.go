package main

import (
	"log"
	"strconv"
)

type Variable struct {
	Value int
}

type Scope struct {
	Variables map[string]Variable
	Parent    *Scope
}

type Function struct {
	Id   string
	Args []string
	Body ExpressionBlock
}

type RetVal struct {
	Value int
}

var functions map[string]Function = make(map[string]Function)
var mainScope Scope = Scope{Variables: make(map[string]Variable), Parent: nil}

func FindVariable(scope *Scope, id string) *Variable {
	v, ok := scope.Variables[id]
	if ok {
		return &v
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

	left := n.Left.Interpret(scope).(int)
	right := n.Right.Interpret(scope).(int)

	switch n.Operation.Type {
	case TOK_PLUS:
		return left + right
	case TOK_MINUS:
		return left - right
	case TOK_MULT:
		return left * right
	case TOK_DIV:
		return left / right
	case TOK_EQUALS:
		if left == right {
			return 1
		}
		return 0
	case TOK_LESS_THAN:
		if left < right {
			return 1
		}
		return 0
	case TOK_GREATER_THAN:
		if left > right {
			return 1
		}
		return 0
	case TOK_AND:
		if left != 0 && right != 0 {
			return 1
		}
		return 0
	case TOK_OR:
		if left != 0 || right != 0 {
			return 1
		}
		return 0
	}
	return nil
}

func (n ExpressionUnaryOp) Interpret(scope *Scope) interface{} {
	if n.Value.Interpret(scope) == nil {
		log.Fatalln("ERROR: Can't operate with a nil type!")
	}
	switch n.Operation.Type {
	case TOK_NOT:
		if n.Value.Interpret(scope).(int) == 0 {
			return 1.0
		}
		return 0.0
	case TOK_MINUS:
		return -(n.Value.Interpret(scope).(int))
	}
	return nil
}

func (n StatementConditional) Interpret(scope *Scope) interface{} {
	switch n.Tok.Type {
	case TOK_IF:
		v := n.Condition.Interpret(scope)
		if v.(int) != 0 {
			return n.Body.Interpret(scope)
		} else {
			if n.Next != nil {
				return n.Next.Interpret(scope)
			}
		}
	case TOK_WHILE:
		for n.Condition.Interpret(scope).(int) != 0 {
			val := n.Body.Interpret(scope)
			if IsType(val, RetVal{}) {
				return val
			}
		}
	case TOK_ELSE:
		return n.Body.Interpret(scope)
	case TOK_ELSIF:
		if n.Condition.Interpret(scope).(int) != 0 {
			return n.Body.Interpret(scope)
		} else {
			return n.Next.Interpret(scope)
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
		scope = &Scope{Variables: make(map[string]Variable), Parent: &mainScope}
	} else {
		scope = &Scope{Variables: make(map[string]Variable), Parent: scope}
	}
	for i := 0; i < len(n.Body); i++ {
		result := n.Body[i].Interpret(scope)
		if result != nil {
			if IsType(result, RetVal{}) {
				return result.(RetVal).Value
			}

			if IsType(n.Body[i], StatementConditional{}) {
				return result
			}
		}
	}
	return nil
}

func (n StatementFunctionDeclaration) Interpret(scope *Scope) interface{} {
	functions[n.Id] = Function{Id: n.Id, Args: n.Args, Body: n.Body}
	return nil
}

func (fun Function) Call(Args []Expression, scope *Scope) int {
	newScope := Scope{Variables: make(map[string]Variable), Parent: &mainScope}

	for i := 0; i < len(fun.Args); i++ {
		newScope.Variables[fun.Args[i]] = Variable{Value: Args[i].Interpret(scope).(int)}
	}

	val := fun.Body.Interpret(&newScope)

	if val == nil {
		log.Printf("WARNING: \"%s\" function should return something!\n", fun.Id)
		return 0
	}

	return val.(int)
}

func (n ExpressionCall) Interpret(scope *Scope) interface{} {
	v, ok := functions[n.Id]
	if !ok {
		log.Fatalln("ERROR: Can't call function", n.Id, "cuz it is not defined!")
	}
	if len(n.Args) != len(v.Args) {
		log.Fatalln("ERROR: Number of arguments is not matching in call to function", n.Id, "!\nExpected arguments count:", len(v.Args), "\nArguments count recived: ", len(n.Args))
	}

	return v.Call(n.Args, scope)
}

func (n ExpressionLiteral) Interpret(scope *Scope) interface{} {
	switch n.Type {
	case num_literal:
		val, err := strconv.Atoi(n.Tok.Lexme)
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

func (n StatementReturn) Interpret(scope *Scope) interface{} {
	return RetVal{Value: n.Value.Interpret(scope).(int)}
}

func (n ExpressionAssigment) Interpret(scope *Scope) interface{} {

	v := n.Value.Interpret(scope)
	if v == nil {
		log.Fatalln("ERROR: Can't set variable", n.Id, "to a nil value!")
	}
	scope.Variables[n.Id] = Variable{Value: v.(int)}
	return scope.Variables[n.Id]
}
