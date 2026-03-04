package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

type (
	SqlDBConfig struct {
		Driver          string `yaml:"driver"`
		Dsn             string `yaml:"dsn"`
		MaxIdleConn     int    `yaml:"maxIdleConn"`
		MaxOpenConn     int    `yaml:"maxOpenConn"`
		MaxConnIdleTime int    `yaml:"maxConnIdleTime"`
		MaxConnLifeTime int    `yaml:"maxConnLifeTime"`
	}
	RedisConfig struct {
		Dsn string `yaml:"dsn"`
	}
	DBConfig struct {
		SqlDBConfig `yaml:"sql"`
		RedisConfig `yaml:"redis"`
	}
)

func initDBConfig() DBConfig {
	return DBConfig{
		SqlDBConfig: SqlDBConfig{
			Dsn:             "",
			MaxIdleConn:     20,
			MaxOpenConn:     0,
			MaxConnIdleTime: 20,
			MaxConnLifeTime: 60,
		},
	}
}

func (db *DBConfig) RedisOptions() *redis.Options {
	ro, err := redis.ParseURL(db.RedisConfig.Dsn)
	if err != nil {
		clog.Warn(fmt.Sprintf("invalid redis dsn: %s", palette.Red(db.RedisConfig.Dsn)))
	}
	return ro
}
