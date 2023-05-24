
[![Go Reference](https://pkg.go.dev/badge/github.com/flachnetz/poly.svg)](https://pkg.go.dev/github.com/flachnetz/poly)

# Poly

Simple polymorphic serialization and deserialization.

```go
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
```
