package util

import (
	"fmt"
	"reflect"
	"sync"
)

type (
	CopierPool sync.Map
	Copier     struct {
		fns []func(dest, src reflect.Value)
	}
)

var copierPool sync.Map

func GetCopier(src, dest interface{}) *Copier {
	cKey := fmt.Sprintf("%v->%v", src, dest)
	if v, ok := copierPool.Load(cKey); ok {
		return v.(*Copier)
	}
	copier := NewCopier(reflect.TypeOf(src), reflect.TypeOf(dest))
	copierPool.Store(cKey, copier)
	return copier
}

// .NewCopier(A,B) -- A,B [struct] <-pointer
func NewCopier(srcType, destType reflect.Type) *Copier {
	var fns []func(dest, src reflect.Value)
	if srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
	}
	if destType.Kind() == reflect.Pointer {
		destType = destType.Elem()
	}
	for i := 0; i < destType.NumField(); i += 1 {
		df := destType.Field(i)
		sf, ok := srcType.FieldByName(df.Name)
		// 未导出字段通过反射赋值行为为panic
		if !ok || !df.IsExported() {
			continue
		}
		if sf.Type.AssignableTo(df.Type) {
			fns = append(fns, func(dest, src reflect.Value) {
				dest.FieldByName(df.Name).Set(src.FieldByName(df.Name))
			})
		} else if sf.Type.ConvertibleTo(df.Type) {
			fns = append(fns, func(dest, src reflect.Value) {
				dest.FieldByName(df.Name).Set(src.FieldByName(df.Name).Convert(df.Type))
			})
		}
	}
	return &Copier{
		fns: fns,
	}
}

// Copy(&B,&A) -- A,B [pointer] <- struct
func (c *Copier) Copy(dest, src interface{}) {
	destValue := reflect.ValueOf(dest)
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Pointer || srcValue.Kind() != reflect.Pointer {
		panic("[copier] - both dest and src must be pointers to struct")
	}

	destValue = destValue.Elem()
	srcValue = srcValue.Elem()

	if destValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		panic("[copier] -  both dest and src must point to structs")
	}

	for _, fn := range c.fns {
		fn(destValue, srcValue)
	}
}
