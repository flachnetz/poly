package poly

import (
	"fmt"
	"reflect"
)

// Type combines a runtime type with a name. This is used to map the the type field on
// a serialized Poly instance to a Type
type Type struct {
	Name        string
	RuntimeType reflect.Type
}

// TypeList provides a list of runtime types.
//
// Implementations of this interface must work directly on an uninitialized type
// if they are supposed to be used with Poly. For example:
//
//	var myTypeList TypeList = MyCustomTypeList{}
//	myTypeList.Types() // must work
//
// The type mapping must be bijective - there must be a one to one mapping between
// Type.Name and Type.RuntimeType
type TypeList interface {
	// Types returns the bijective type mappings
	Types() []Type
}

// TypeItem is a compile time generic TypeList item. It contains the first type
// of a TypeList as well as the remaining types (Next).
// See https://en.wikipedia.org/wiki/Cons for this kind of recursive list definition.
type TypeItem[T any, Next TypeList] struct{}

func (c TypeItem[T, Next]) Types() []Type {
	var next Next
	return append([]Type{typeOf[T]()}, next.Types()...)
}

// Nil is a TypeList of length zero.
type Nil struct{}

func (n Nil) Types() []Type {
	return nil
}

// typeOf returns the Type of the type parameter T.
func typeOf[T any]() Type {
	var value T
	runtimeType := reflect.TypeOf(&value).Elem()

	return Type{
		Name:        typeNameOf(runtimeType),
		RuntimeType: runtimeType,
	}
}

func typeNameOf(typ reflect.Type) string {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
}
