package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/wendisx/puzzle/pkg/clog"
)

// 并不是所有的database驱动支持从LastInsertId中得到合适的pkey值, 这里以最宽松的方式执行插入
func InsertWithPlace(ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

func UpdateWithPlace(ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

func DeleteWithPlace(ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

func QueryWithPlace[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) (R, error) {
	var dest R
	err := db.GetContext(ctx, &dest, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

func QListWithPlace[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) ([]R, error) {
	var dest []R
	err := db.SelectContext(ctx, &dest, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

func QPageWithPlace[R any](ctx context.Context, db *sqlx.DB, sqlStr string, args ...any) (Page[R], error) {
	var page Page[R]
	list := make([]R, 0)
	list, err := QListWithPlace[R](ctx, db, sqlStr, args...)
	if err != nil {
		clog.Error(err.Error())
	}
	// 传递时, 大概率offset和count是args[len(args)-2]和args[len(args)-1], 这里不尝试赋值page, 依赖上层行为处理.
	page.Items = list
	return page, err
}

func InsertWithName(ctx context.Context, db *sqlx.DB, sqlStr string, obj any) error {
	_, err := db.NamedExecContext(ctx, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

func UpdateWithName(ctx context.Context, db *sqlx.DB, sqlStr string, obj any) error {
	_, err := db.NamedExecContext(ctx, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

func DeleteWithName(ctx context.Context, db *sqlx.DB, sqlStr string, obj any) error {
	_, err := db.NamedExecContext(ctx, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
		return err
	}
	return nil
}

func QueryWithName[R any](ctx context.Context, db *sqlx.DB, sqlStr string, obj any) (R, error) {
	var dest R
	rows, err := db.NamedQueryContext(ctx, sqlStr, obj)
	if err == nil {
		for rows.Next() {
			err = rows.StructScan(&dest)
		}
	}
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

func QListWithName[R any](ctx context.Context, db *sqlx.DB, sqlStr string, obj any) ([]R, error) {
	dest := make([]R, 0)
	rows, err := db.NamedQueryContext(ctx, sqlStr, obj)
	var row R
	for rows.Next() {
		err = rows.StructScan(&row)
		if err != nil {
			clog.Error(err.Error())
			continue
		}
		dest = append(dest, row)
	}
	if err != nil {
		clog.Error(err.Error())
	}
	return dest, err
}

func QPageWithName[R any](ctx context.Context, db *sqlx.DB, sqlStr string, obj any) (Page[R], error) {
	var page Page[R]
	list := make([]R, 0)
	list, err := QListWithName[R](ctx, db, sqlStr, obj)
	if err != nil {
		clog.Error(err.Error())
	}
	page.Items = list
	return page, err
}
