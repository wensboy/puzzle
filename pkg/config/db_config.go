package config

type (
	SqlDBConfig struct {
		Driver          string `yaml:"driver"`
		Dsn             string `yaml:"dsn"`
		MaxIdleConn     int    `yaml:"maxIdleConn"`
		MaxOpenConn     int    `yaml:"maxOpenConn"`
		MaxConnIdleTime int    `yaml:"maxConnIdleTime"`
		MaxConnLifeTime int    `yaml:"maxConnLifeTime"`
	}
	DBConfig struct {
		SqlDBConfig `yaml:"sql"`
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
