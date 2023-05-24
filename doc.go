/*
Package poly provides polymorphic serialization and deserialization to and from json for
go values. It archives this by serializing a value that is wrapped in a Poly container together
with the values type name, for example:

	{"type": "mypackage.Dog", "value": {"name": "Wuffy", "barks": true}}

A little bit of work is required for golang to access the possible types for deserialization into an interface.
This library provides the TypeList interface together with the generic TypeItem and Nil types, to build up
a variable length list of possible deserialization types during compile time.

Given struct Cat and Dog, both of interface type Animal, a Poly type can be constructed:

	type MyState struct {
		Animal poly.Poly[Animal, TypeItem[Dog, TypeItem[Cat, Nil]]]
	}

Setting the Poly.Value field to a Dog or Cat instance will serialize the value and deserialize the
example above back into a Cat or Dog value in the Poly.Value field.

A full example:

	type Speaker interface{ Speak() }

	type Dog struct { IsDog bool }
	func (Dog) Speak() {}

	type Cat struct { IsCat bool }
	func (Cat) Speak() {}

	var animal poly.Poly[Speaker, poly.TypeList[Dog, poly.TypeList[Cat, poly.Nil]]]

	// write a dog to json
	animal.Value = Dog{}
	bytes, _ := json.Marshal(animal)

	// deserialize the dog
	_ = json.Unmarshal(bytes, &animal)
	animal.Value.Speak() // wuff
*/
package poly
