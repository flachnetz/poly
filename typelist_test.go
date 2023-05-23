package poly

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type (
	A struct{}
	B struct{}
	I interface{}
)

var _ = Describe("TypeOf", func() {
	It("should work for structs", func() {
		tA := reflect.TypeOf((*A)(nil)).Elem()
		tB := reflect.TypeOf((*B)(nil)).Elem()

		Expect(typeOf[A]()).To(Equal(tA))
		Expect(typeOf[B]()).To(Equal(tB))
	})

	It("should work for interfaces", func() {
		tI := reflect.TypeOf((*I)(nil)).Elem()
		Expect(typeOf[I]()).To(Equal(tI))
	})

	It("should work for pointers", func() {
		tA := reflect.TypeOf((**A)(nil)).Elem()
		tI := reflect.TypeOf((**I)(nil)).Elem()
		Expect(typeOf[*A]()).To(Equal(tA))
		Expect(typeOf[*I]()).To(Equal(tI))
	})
})

var _ = Describe("Typelist", func() {
	It("should return no items for an empty list", func() {
		var t TypeList = Nil{}
		Expect(t.Types()).To(BeEmpty())
	})

	It("should return the types in the correct order", func() {
		var t TypeItem[A, TypeItem[B, TypeItem[I, Nil]]]

		Expect(t.Types()).To(
			Equal([]reflect.Type{
				reflect.TypeOf((*A)(nil)).Elem(),
				reflect.TypeOf((*B)(nil)).Elem(),
				reflect.TypeOf((*I)(nil)).Elem(),
			}),
		)
	})
})
