package config

type (
	DBConfig struct {
		Dsn             string `yaml:"dsn"`
		MaxIdleConn     int    `yaml:"maxIdleConn"`
		MaxOpenConn     int    `yaml:"maxOpenConn"`
		MaxConnIdleTime int    `yaml:"maxConnIdleTime"`
		MaxConnLifeTime int    `yaml:"maxConnLifeTime"`
	}
	dbConfigOption func(c *DBConfig)
)

func initDBConfig(dsn string) DBConfig {
	return DBConfig{
		Dsn:             dsn,
		MaxIdleConn:     20,
		MaxOpenConn:     0,
		MaxConnIdleTime: 20,
		MaxConnLifeTime: 60,
	}
}

func (c *DBConfig) SetupConfig(opts ...dbConfigOption) {
	for _, fn := range opts {
		fn(c)
	}
}
