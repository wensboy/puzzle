package config

import (
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

	DICTKEY_CONFIG  = "_dict_config_"
	DICTKEY_COMMAND = "_dict_command"

	DATAKEY_CONFIG = "_data_config_"
	DATAKEY_ENV    = "_data_env_"
	DATAKEY_CLI    = "_data_cli_"
)

var (
	_dict_directory *DictDirectory
	_dict_timeout   = 10 * time.Minute
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
	dd, found := _dict_directory.find(string(k))
	if !found {
		clog.Panic(fmt.Sprintf("from dict directory can't find dict_key(%s)", palette.Red(k)))
	}
	return dd
}

func HasDict(k DictKey) bool {
	_, found := _dict_directory.find(string(k))
	return found
}

// Create the broadest data dictionary and only allow this data dictionary to enter _dict_directory.
func (ds *DictDirectory) record(k string, dict DataDict[any]) {
	ds.qdict.Store(k, dict)
	clog.Info(fmt.Sprintf("put dict_key(%s) into dict directory", palette.SkyBlue(k)))
}

func (ds *DictDirectory) find(k string) (DataDict[any], bool) {
	if v, ok := ds.qdict.Load(k); ok {
		return v.(DataDict[any]), true
	}
	return DataDict[any]{}, false
}

// next dict ttl config
func NextDictTTL(ttl time.Duration) {
	_dict_timeout = ttl
}

// common usage data dict
func NewDataDict[V any](name DictKey) DataDict[V] {
	dd := ttlcache.New(
		ttlcache.WithCapacity[string, V](_dict_capacity),
		ttlcache.WithTTL[string, V](_dict_timeout),
	)
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
	clog.Info(fmt.Sprintf("put data_key(%s) into dict_key(%s)", palette.SkyBlue(k), palette.SkyBlue(dd.name)))
}

func (dd *DataDict[V]) Find(k string) *ttlcache.Item[string, V] {
	if !dd.dict.Has(k) {
		clog.Panic(fmt.Sprintf("from dict_key(%s) can't find data_key(%s)", palette.Red(dd.name), palette.Red(k)))
	}
	return dd.dict.Get(k)
}

func (dd *DataDict[V]) Remove(k string) {
	dd.dict.Delete(k)
}

func (dd *DataDict[V]) RemoveAll() {
	dd.dict.DeleteAll()
}
