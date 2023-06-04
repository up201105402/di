package steps

import (
	"fmt"
	"reflect"
)

var stepTypeRegistry = make(map[string]reflect.Type)

func init() {
	stepTypes := []interface{}{CheckoutRepoStep{}}
	for _, v := range stepTypes {
		stepTypeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}
}

func makeInstance(name string) interface{} {
	v := reflect.New(stepTypeRegistry[name]).Elem()
	// Maybe fill in fields here if necessary
	return v.Interface()
}
