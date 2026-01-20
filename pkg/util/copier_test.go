package util

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/wendisx/puzzle/pkg/clog"
)

// [passed]
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

// [passed]
func Test_basic_type(t *testing.T) {
	var a int64 = 1
	var b float64
	cp := GetCopier(a, b)
	cp.Copy(&b, &a)
	clog.Info(fmt.Sprintf("%f", b))
}

// []
func Test_uuid(t *testing.T) {
	type (
		A struct {
			Id uuid.UUID
		}
		B struct {
			Id uuid.UUID
		}
	)
	a := A{
		Id: uuid.New(),
	}
	var b B
	cp := GetCopier(a, b)
	cp.Copy(&b, &a)
	clog.Info(fmt.Sprintf("a = %+v", a))
	clog.Info(fmt.Sprintf("b = %+v", b))
}

// [passed]
func Test_copy_slice(t *testing.T) {
	type (
		A struct {
			Id   uint64
			Name string
		}
		ADao struct {
			Id   uint64
			Name string
			Age  int
		}
	)
	LA := []A{}
	for i := 0; i < 10; i += 1 {
		LA = append(LA, A{Id: uint64(i), Name: fmt.Sprintf("A_%d", i)})
	}
	clog.Info(fmt.Sprintf("%#v", LA))
	clog.Info(fmt.Sprintf("%#v", CopySlice[A, ADao](LA)))
}
