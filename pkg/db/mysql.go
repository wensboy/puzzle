package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/wendisx/puzzle/pkg/errors"
)

var mysqlDB *sqlx.DB

const (
	codeSqlExecFail = 50001
	codeUnexpected  = 50002

	msgSqlExecFail   = "sql execute fail"
	msgSqlUnexpected = "unexpected sql meta"
)

const (
	driverMysql = "mysql"
	// dsnTemp     = "[user:password]@porto(host:port)[/db?options...]"
	defaultDsn = "root:root@tcp(127.0.0.1:3306)/main?charset=utf8mb4&parseTime=true&loc=Local&timeout=3s"

	ConnTimeout = 3 * time.Second
)

func InitMysqlDB(dsn string) error {
	if dsn == "" {
		dsn = defaultDsn
	}
	ctx, cancle := context.WithTimeout(context.Background(), ConnTimeout)
	defer cancle()
	var err error
	mysqlDB, err = sqlx.ConnectContext(ctx, driverMysql, dsn)
	if err != nil {
		return err
	}
	return nil
}

func GetMysqlDB() *sqlx.DB {
	return mysqlDB
}

// A standard MySQL database operation should at least implement minimal operations,
// always perform database operations in a combined manner, and use transaction control
// behavior when necessary.
// Basic operations =>
// InsertOne(ctx, obj) (id,error)
// UpdateOne(ctx, obj) (id,error)
// DeleteOne(ctx, obj) (error)
// QueryOneById(ctx, id) (obj, error)
// QueryOneByxxx(ctx, xxx) (obj, error)
// tx may start =>
// InsertList(ctx, []obj) ([]id,error)
// DeleteList(ctx, []id) error
// QueryList(ctx, ...) ([]obj,error)
// QueryPage(ctx, *page[obj]) (error)

func InsertOne[T any](ctx context.Context, db *sqlx.DB, em uint8, sqlStr string, obj T, args ...any) (int64, error) {
	var res sql.Result
	var err error
	switch em {
	case Place:
		res, err = db.ExecContext(ctx, sqlStr, args...)
	case Named:
		res, err = db.NamedExecContext(ctx, sqlStr, obj)
	}
	if err != nil {
		return 0, errors.NewIErr(codeSqlExecFail, msgSqlExecFail)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.NewIErr(codeUnexpected, msgSqlUnexpected)
	}
	return id, nil
}

func UpdateOne[T any](ctx context.Context, db *sqlx.DB, em uint8, sqlStr string, obj T, args ...any) error {
	var res sql.Result
	var err error
	switch em {
	case Place:
		res, err = db.ExecContext(ctx, sqlStr, args...)
	case Named:
		res, err = db.NamedExecContext(ctx, sqlStr, obj)
	}
	if err != nil {
		return errors.NewIErr(codeSqlExecFail, msgSqlExecFail)
	}
	cnt, err := res.RowsAffected()
	if err != nil || (cnt != 0 && cnt != 1) {
		return errors.NewIErr(codeUnexpected, msgSqlUnexpected)
	}
	return nil
}

func DeleteOne[T any](ctx context.Context, db *sqlx.DB, em uint8, sqlStr string, obj T, args ...any) error {
	var res sql.Result
	var err error
	switch em {
	case Place:
		res, err = db.ExecContext(ctx, sqlStr, args...)
	case Named:
		res, err = db.NamedExecContext(ctx, sqlStr, obj)
	}
	if err != nil {
		return errors.NewIErr(codeSqlExecFail, msgSqlExecFail)
	}
	cnt, err := res.RowsAffected()
	if err != nil || (cnt != 0 && cnt != 1) {
		return errors.NewIErr(codeUnexpected, msgSqlUnexpected)
	}
	return nil
}

func QueryOne[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) (R, error) {
	var dest R
	var err error
	err = db.GetContext(ctx, &dest, sqlStr, args...)
	if err != nil {
		return dest, errors.NewIErr(codeSqlExecFail, msgSqlExecFail)
	}
	return dest, nil
}

func QueryList[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) ([]R, error) {
	var dest []R
	var err error
	err = db.SelectContext(ctx, &dest, sqlStr, args...)
	if err != nil {
		return dest, errors.NewIErr(codeSqlExecFail, msgSqlExecFail)
	}
	return dest, nil
}

func QueryPage[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) (Page[R], error) {
	var dest Page[R]
	var err error
	err = db.SelectContext(ctx, &dest.Items, sqlStr, args...)
	if err != nil {
		return dest, errors.NewIErr(codeSqlExecFail, msgSqlExecFail)
	}
	return dest, nil
}
