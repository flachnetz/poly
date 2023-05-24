package poly

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Poly is a transparent wrapper around a generic value of type I that uses a TypeList
// to serialize and deserialize different subtypes of I. I should be an interface.
// If Value is nil, the container serializes to json 'null'.
type Poly[I any, Tl TypeList] struct {
	Value I
}

type envelope[I any] struct {
	Type  string `json:"type"`
	Value I      `json:"value"`
}

func (c *Poly[I, Tl]) MarshalJSON() ([]byte, error) {
	rValue := reflect.ValueOf(c.Value)
	if !rValue.IsValid() {
		return []byte("null"), nil
	}

	var typeList Tl
	for _, typ := range typeList.Types() {
		if typ.RuntimeType == rValue.Type() {
			return json.Marshal(envelope[I]{
				Type:  typ.Name,
				Value: c.Value,
			})
		}
	}

	return nil, fmt.Errorf("no type mapping in TypeList found for %s", rValue.Type().String())
}

func (c *Poly[I, Tl]) UnmarshalJSON(bytes []byte) error {
	// json 'null' is deserialized to a nil Value of T
	if len(bytes) == 4 && string(bytes) == "null" {
		var nilValue I
		c.Value = nilValue
		return nil
	}
	// read into envelope to get the type
	var envelope envelope[json.RawMessage]
	if err := json.Unmarshal(bytes, &envelope); err != nil {
		return err
	}

	// get the possible target types from the type list
	var typeList Tl
	types := typeList.Types()

	for _, typ := range types {
		if typ.Name != envelope.Type {
			continue
		}

		// found our target type
		ptrInstance := reflect.New(typ.RuntimeType)

		// deserialize to value
		if err := json.Unmarshal(envelope.Value, ptrInstance.Interface()); err != nil {
			return err
		}

		// convert to the target type
		value, ok := ptrInstance.Elem().Interface().(I)
		if !ok {
			return fmt.Errorf("%q does not implement the required interface", typ.Name)
		}

		c.Value = value

		return nil
	}

	return fmt.Errorf("did not find a type for %q", envelope.Type)
}
