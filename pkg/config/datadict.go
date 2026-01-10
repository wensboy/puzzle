package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

/*
	DataDict -- a simple internal data collection
	1. Allows the creation of custom data dictionaries, but this is just a wrapper for ttlcache.
	2. An extremely simple dictionary processing design, relying on generics and assertion design.
*/

const (
	// some default key from internal, exactly I'd like to got them from some .json files. :)
	_dict_capacity = 1 << 10
	_dict_timeout  = 10 * time.Minute

	DICTKEY_CONFIG DictKey = "_dict_config_"
)

var (
	_dict_directory *DictDirectory
)

type (
	DictKey       string
	DictDirectory struct {
		qdict sync.Map
	}
	DataDict[V any] struct {
		name string
		dict *ttlcache.Cache[string, V]
	}
)

func init() {
	_dict_directory = new(DictDirectory)
}

func PutDict(k DictKey, dict DataDict[any]) {
	_dict_directory.record(string(k), dict)
}

// panic from  .get(k)
func GetDict(k DictKey) DataDict[any] {
	return _dict_directory.find(string(k))
}

// Create the broadest data dictionary and only allow this data dictionary to enter _dict_directory.
func (ds *DictDirectory) record(k string, dict DataDict[any]) {
	ds.qdict.Store(k, dict)
}

func (ds *DictDirectory) find(k string) DataDict[any] {
	if v, ok := ds.qdict.Load(k); ok {
		return v.(DataDict[any])
	}
	clog.Panic(fmt.Sprintf("try to get invalid datadict from dict queue with dict_key(%s)", palette.Red(k)))
	return DataDict[any]{
		name: "unreachable_code",
	}
}

// common usage data dict
func NewDataDict[V any](name DictKey) DataDict[V] {
	dd := ttlcache.New(
		ttlcache.WithCapacity[string, V](_dict_capacity),
		ttlcache.WithTTL[string, V](_dict_timeout),
	)
	dd.OnInsertion(func(ctx context.Context, i *ttlcache.Item[string, V]) {
		clog.Info(fmt.Sprintf("DataDict(%s) on Insert with [%s, %v]", palette.SkyBlue(name), i.Key(), i.Value()))
	})
	dd.OnUpdate(func(ctx context.Context, i *ttlcache.Item[string, V]) {
		clog.Info(fmt.Sprintf("DataDict(%s) on Update with [%s, %v]", palette.SkyBlue(name), i.Key(), i.Value()))
	})
	dd.OnEviction(func(ctx context.Context, er ttlcache.EvictionReason, i *ttlcache.Item[string, V]) {
		clog.Info(fmt.Sprintf("DataDict(%s) on Eviction with [%s, %v]", palette.SkyBlue(name), i.Key(), i.Value()))
	})
	return DataDict[V]{
		name: string(name),
		dict: dd,
	}
}

func (dd *DataDict[V]) Name() DictKey {
	return DictKey(dd.name)
}

func (dd *DataDict[V]) Record(k string, v V) {
	dd.dict.Set(k, v, ttlcache.DefaultTTL)
}

func (dd *DataDict[V]) Find(k string) *ttlcache.Item[string, V] {
	return dd.dict.Get(k)
}

func (dd *DataDict[V]) Remove(k string) {
	dd.dict.Delete(k)
}

func (dd *DataDict[V]) RemoveAll() {
	dd.dict.DeleteAll()
}
