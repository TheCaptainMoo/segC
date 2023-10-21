package main

import (
	"reflect"
)

func IsType(variable interface{}, target any) bool {
	varType := reflect.TypeOf(variable)
	return varType == reflect.TypeOf(target)
}
