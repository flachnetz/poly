package poly

import "reflect"

type TypeList interface {
	Types() []reflect.Type
}

// TypeItem is a generic TypeList item. It contains the first type
// of a TypeList as well as the remaining types (Next) during compile time.
// See https://www.wikiwand.com/en/Cons for this kind of recursive list definition.
type TypeItem[T any, Next TypeList] struct{}

func (c TypeItem[T, Next]) Types() []reflect.Type {
	var next Next

	return append([]reflect.Type{typeOf[T]()}, next.Types()...)
}

// Nil is a TypeList of length zero.
type Nil struct{}

func (n Nil) Types() []reflect.Type {
	return nil
}

// typeOf returns the reflect.Type of the type parameter T.
func typeOf[T any]() reflect.Type {
	var value T
	return reflect.TypeOf(&value).Elem()
}
