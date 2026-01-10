package database

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wendisx/puzzle/pkg/clog"
)

/*
mysql database
*/
const (
	_driver_mysql = "mysql"
	// _dsn_template = "<user>:<password>@<proto>(<host>:<port>)[/<db_name>][?options]"
	_default_dsn  = "root:root@tcp(127.0.0.1:3306)/test?useSSL=false&loc=local" // 最终假设运行的mysql实例
	_conn_timeout = 3 * time.Second
)

type (
	MysqlDB     *sqlx.DB
	MysqlOption func(MysqlDB)
)

func InitMysql(dsn string) MysqlDB {
	if dsn == "" {
		dsn = _default_dsn
	}
	ctx, cancle := context.WithTimeout(context.Background(), _conn_timeout)
	defer cancle()
	db, err := sqlx.ConnectContext(ctx, _driver_mysql, dsn)
	if err != nil {
		clog.Panic(err.Error())
	} else {
		clog.Info("<pkg.database> Mysql Database initialization successful.")
	}
	return db
}

func SetupMysql(mdb MysqlDB, opts ...MysqlOption) {
	for _, fn := range opts {
		fn(mdb)
	}
}
