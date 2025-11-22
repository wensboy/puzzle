package util

import (
	"log"
	"testing"
)

func TestCopier(t *testing.T) {
	type (
		C struct {
			id   int
			Name string
		}
		D struct {
			id   int
			Name string
		}
		A struct {
			C
			Eid   []byte
			Email string
		}
		B struct {
			C
			Eid   []byte
			Email string
		}
	)
	t.Run("exported and em", func(t *testing.T) {
		a := A{
			C: C{
				id:   1,
				Name: "a",
			},
			Eid:   []byte{'t', 'h', 'e', ' ', 'A'},
			Email: "A@xx.com",
		}
		b := B{}
		GetCopier(A{}, B{}).Copy(&b, &a)
		log.Printf("%#v\n", b)
	})
	t.Run("unexported", func(t *testing.T) {
		c := C{
			id:   3,
			Name: "c",
		}
		d := D{
			id: 4,
		}
		GetCopier(C{}, D{}).Copy(&d, &c)
		log.Printf("%#v\n", d)
	})
}
