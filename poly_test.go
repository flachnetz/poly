package poly

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Speaker interface{ __sealed() }

type Dog struct {
	Speaker `json:"-"`
	IsDog   bool
}

type Cat struct {
	Speaker `json:"-"`
	IsCat   bool
}

var _ = Describe("Poly", func() {
	It("can serialize", func() {
		var c Poly[Speaker, TypeItem[Dog, TypeItem[Cat, Nil]]]
		c.Value = Dog{IsDog: true}

		bytes, err := c.MarshalJSON()
		Expect(bytes, err).To(MatchJSON(`{"type": "github.com/flachnetz/poly.Dog", "value": {"IsDog": true}}`))
	})

	It("can serialize a ptr value", func() {
		var c Poly[Speaker, TypeItem[*Dog, TypeItem[*Cat, Nil]]]
		c.Value = &Dog{IsDog: true}

		bytes, err := c.MarshalJSON()
		Expect(bytes, err).To(MatchJSON(`{"type": "github.com/flachnetz/poly.Dog", "value": {"IsDog": true}}`))
	})

	It("can serialize nil value", func() {
		var c Poly[Speaker, TypeItem[Dog, TypeItem[Cat, Nil]]]
		bytes, err := c.MarshalJSON()
		Expect(bytes, err).To(MatchJSON(`null`))
	})

	It("can serialize nil ptr value", func() {
		var c Poly[Speaker, TypeItem[*Dog, Nil]]
		c.Value = (*Dog)(nil)
		bytes, err := c.MarshalJSON()
		Expect(bytes, err).To(MatchJSON(`{"type": "github.com/flachnetz/poly.Dog", "value": null}`))
	})

	It("can deserialize 'null' value for non-ptr type", func() {
		var c Poly[Speaker, TypeItem[Dog, TypeItem[Cat, Nil]]]
		err := c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Dog", "value": null}`))
		Expect(c.Value, err).To(Equal(Dog{}))
	})

	It("can deserialize 'null' value for ptr type", func() {
		var c Poly[Speaker, TypeItem[*Dog, TypeItem[Cat, Nil]]]
		err := c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Dog", "value": null}`))
		Expect(c.Value, err).To(Equal((*Dog)(nil)))
	})

	It("can deserialize 'null' value", func() {
		var c Poly[Speaker, TypeItem[Dog, TypeItem[Cat, Nil]]]
		err := c.UnmarshalJSON([]byte(`null`))
		Expect(c.Value, err).To(BeNil())
	})

	It("can deserialize to the right type", func() {
		var c Poly[Speaker, TypeItem[Dog, TypeItem[Cat, Nil]]]
		err := c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Dog", "value": {"IsDog": true}}`))
		Expect(c.Value, err).To(Equal(Dog{IsDog: true}))

		err = c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Cat", "value": {"IsCat": true}}`))
		Expect(c.Value, err).To(Equal(Cat{IsCat: true}))
	})

	It("can deserialize a ptr value", func() {
		var c Poly[Speaker, TypeItem[*Dog, Nil]]
		err := c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Dog", "value": {"IsDog": true}}`))
		Expect(c.Value, err).To(Equal(&Dog{IsDog: true}))
	})

	It("fail if the type doesnt match the interface", func() {
		// not a Speaker
		type Fish struct{}

		var c Poly[Speaker, TypeItem[Fish, Nil]]
		err := c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Fish", "value": {}}`))
		Expect(err).To(HaveOccurred())
	})

	It("works with any interface", func() {
		var c Poly[any, TypeItem[int, TypeItem[string, Nil]]]
		err := c.UnmarshalJSON([]byte(`{"type": ".int", "value": 2}`))
		Expect(c.Value, err).To(Equal(2))

		err = c.UnmarshalJSON([]byte(`{"type": ".string", "value": "foo"}`))
		Expect(c.Value, err).To(Equal("foo"))
	})

	It("fails if type is not known", func() {
		var c Poly[Speaker, TypeItem[Dog, TypeItem[Cat, Nil]]]
		err := c.UnmarshalJSON([]byte(`{"type": "github.com/flachnetz/poly.Mouse", "value": {}}`))
		Expect(err).To(HaveOccurred())
	})
})
