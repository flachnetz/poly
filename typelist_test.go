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
		tA := Type{
			RuntimeType: reflect.TypeOf((*A)(nil)).Elem(),
			Name:        "github.com/flachnetz/poly.A",
		}
		tB := Type{
			RuntimeType: reflect.TypeOf((*B)(nil)).Elem(),
			Name:        "github.com/flachnetz/poly.B",
		}

		Expect(typeOf[A]()).To(Equal(tA))
		Expect(typeOf[B]()).To(Equal(tB))
	})

	It("should work for interfaces", func() {
		tI := Type{
			RuntimeType: reflect.TypeOf((*I)(nil)).Elem(),
			Name:        "github.com/flachnetz/poly.I",
		}
		Expect(typeOf[I]()).To(Equal(tI))
	})

	It("should work for pointers", func() {
		tA := Type{
			RuntimeType: reflect.TypeOf((**A)(nil)).Elem(),
			Name:        "github.com/flachnetz/poly.A",
		}
		tI := Type{
			RuntimeType: reflect.TypeOf((**I)(nil)).Elem(),
			Name:        "github.com/flachnetz/poly.I",
		}
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
			Equal([]Type{
				{RuntimeType: reflect.TypeOf((*A)(nil)).Elem(), Name: "github.com/flachnetz/poly.A"},
				{RuntimeType: reflect.TypeOf((*B)(nil)).Elem(), Name: "github.com/flachnetz/poly.B"},
				{RuntimeType: reflect.TypeOf((*I)(nil)).Elem(), Name: "github.com/flachnetz/poly.I"},
			}),
		)
	})
})
