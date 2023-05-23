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

func (c *Poly[I, Tl]) MarshalJSON() ([]byte, error) {
	type Envelope struct {
		Typ   string `json:"type"`
		Value I      `json:"value"`
	}

	rValue := reflect.ValueOf(c.Value)
	if !rValue.IsValid() {
		return []byte("null"), nil
	}

	return json.Marshal(Envelope{
		Typ:   typeNameOf(rValue.Type()),
		Value: c.Value,
	})
}

func (c *Poly[I, Tl]) UnmarshalJSON(bytes []byte) error {
	// json 'null' is deserialized to a nil Value of T
	if len(bytes) == 4 && string(bytes) == "null" {
		var nilValue I
		c.Value = nilValue
		return nil
	}

	var envelope struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	// read into envelope to get the type
	if err := json.Unmarshal(bytes, &envelope); err != nil {
		return err
	}

	// get the possible target types from the type list
	var typeList Tl
	types := typeList.Types()

	for _, typ := range types {
		if typeNameOf(typ) != envelope.Type {
			continue
		}

		// found our target type
		var ptrInstance = reflect.New(typ)

		// deserialize to value
		if err := json.Unmarshal(envelope.Value, ptrInstance.Interface()); err != nil {
			return err
		}

		// convert to the target type
		value, ok := ptrInstance.Elem().Interface().(I)
		if !ok {
			return fmt.Errorf("%s does not implement the required interface", typ.String())
		}

		c.Value = value

		return nil
	}

	return fmt.Errorf("did not find a type for %q", envelope.Type)
}

func typeNameOf(typ reflect.Type) string {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
}
