package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/config"
	"github.com/wendisx/puzzle/pkg/palette"
)

func InitRedisDB(ro *redis.Options) *redis.Client {
	if ro == nil {
		ro = config.GetConfig().DBConfig.RedisOptions()
	}
	if ro == nil {
		clog.Error("init redis database fail for invalid redis options")
	}
	rc := redis.NewClient(ro)
	clog.Info("init redis database")
	return rc
}

func GetRedisDB() *redis.Client {
	rdb := InitRedisDB(nil)
	if config.HasDict(config.DICTKEY_CLIENT) {
		clientDict := config.GetDict(config.DICTKEY_CLIENT)
		if clientDict.Has(config.DATAKEY_DB_REDIS) {
			rdb = clientDict.Find(config.DATAKEY_DB_REDIS).Value().(*redis.Client)
		} else {
			clientDict.Record(config.DATAKEY_DB_REDIS, rdb)
		}
	} else {
		clog.Warn(fmt.Sprintf("not exists dict_key(%s) to store data_key(%s)", palette.Red(config.DICTKEY_CLIENT), palette.Red(config.DATAKEY_DB_REDIS)))
	}
	return rdb
}
