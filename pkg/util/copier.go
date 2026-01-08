package util

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/wendisx/puzzle/pkg/clog"
)

const (
	_default_ttl      = 30 * time.Minute
	_default_capacity = 1 << 7
)

// key(string): src_to_dest value(*Copier): pot some functions collection
var copierPool *ttlcache.Cache[string, *Copier]

type (
	Copier struct {
		fns []func(dest, src reflect.Value)
	}
)

func init() {
	copierPool = ttlcache.New(
		ttlcache.WithTTL[string, *Copier](_default_ttl),
		ttlcache.WithCapacity[string, *Copier](_default_capacity),
	)
}

func GetCopier(src, dest interface{}) *Copier {
	cKey := fmt.Sprintf("%v_to_%v", src, dest)
	if copierPool.Has(cKey) {
		return copierPool.Get(cKey).Value()
	}
	copier := NewCopier(reflect.TypeOf(src), reflect.TypeOf(dest))
	copierPool.Set(cKey, copier, ttlcache.DefaultTTL)
	return copier
}

func NewCopier(srcType, destType reflect.Type) *Copier {
	var fns []func(dest, src reflect.Value)
	if srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
	}
	if destType.Kind() == reflect.Pointer {
		destType = destType.Elem()
	}
	// 对于dest中的字段尽可能在src中找到源
	for i := 0; i < destType.NumField(); i += 1 {
		df := destType.Field(i)
		sf, ok := srcType.FieldByName(df.Name)
		// 没找到跳过字段
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

func (c *Copier) Copy(dest, src interface{}) {
	destValue := reflect.ValueOf(dest)
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Pointer || srcValue.Kind() != reflect.Pointer {
		clog.Panic("[copier] - both dest and src must be pointers to struct")
	}

	destValue = destValue.Elem()
	srcValue = srcValue.Elem()

	if destValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		clog.Panic("[copier] -  both dest and src must point to structs")
	}

	for _, fn := range c.fns {
		fn(destValue, srcValue)
	}
}
