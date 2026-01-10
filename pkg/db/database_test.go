package database

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	_test_default_dsn     = "root:ff5740f5-5083-4f71-9ac7-c13144c87b78@tcp(127.0.0.1:53306)/test?charset=utf8mb4"
	_test_ctx_tomeout     = 10 * time.Second
	_test_user_name_basic = "test_user_"
	_test_user_max_limit  = 10
	_test_user_id_basic   = 1
)

var (
	_test_mysql_db MysqlDB
	_test_sql_list = []string{
		`
		insert into user_basic(extern_id,user_name,user_password)
		values 
		(uuid_to_bin(uuid()),?,sha2(?, 256))
		`,
		`
		insert into user_basic(extern_id,user_name,user_password)
		values 
		(uuid_to_bin(uuid()),:user_name,sha2(user_password, 256))
		`,
		`
		insert into user_detail(user_id,nickname,phone,email,avatar)
		values 
		(?,?,?,?,?)
		`,
		`
		insert into user_detail(user_id,nickname,phone,email,avatar)
		values 
		(:user_id,:nickname,:phone,:email,:avatar)
		`,
		`
		update user_basic
		set user_password = sha(?, 256)
		where user_name = ? and deleted = 0
		`,
		`
		update user_basic
		set user_password = sha(:user_password, 256)
		where user_name = :user_name and deleted = 0
		`,
		`
		update user_detail
		set nickname = ?
		where user_id = ?
		`,
		`
		update user_basic
		set deleted = 1
		where user_name = ? and deleted = 0
		`,
		`
		delete from user_basic where deleted = 1 cascade
		`,
		`
		select user_name, user_password 
		from user_basic
		where user_name = ? and deleted = 0
		`,
		`
		select user_id, nickname, phone, email, avatar
		from user_detail
		where user_id = :user_id
		`,
		`
		select user_name, user_password
		from user_basic
		where user_name like ? and deleted = 0
		`,
		`
		select user_id, nickname phone, email, avatar
		from user_detail
		where user_id > :min_user_id
		`,
		`
		select user_id, nickname, phone, email, avatar
		from user_detail
		order by user_id desc
		limit ?,?
		`,
	}
)

type (
	UserBasic struct {
		MysqlMeta
		UserName     string `db:"user_name"`
		UserPassword string `db:"user_password"`
	}
	UserDetail struct {
		MysqlMeta
		UserId   uint64 `db:"user_id"`
		Nickname string `db:"nickname"`
		Phone    string `db:"phone"`
		Email    string `db:"email"`
		Avatar   string `db:"avatar"`
	}
)

// pass
func Test_connect_db(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
}

// pass
func Test_insert_one(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
	ctx, cancle := context.WithTimeout(context.Background(), _test_ctx_tomeout)
	defer cancle()
	// for user_basic
	for i := 0; i < _test_user_max_limit; i += 1 {
		newUser := UserBasic{
			UserName:     fmt.Sprintf("%s%d", _test_user_name_basic, i),
			UserPassword: fmt.Sprintf("%s%d", _test_user_name_basic, i),
		}
		var err error
		if i%2 == 0 {
			err = InsertWithPlace(ctx, _test_mysql_db, _test_sql_list[0], newUser.UserName, newUser.UserPassword)
		} else {
			err = InsertWithName(ctx, _test_mysql_db, _test_sql_list[1], newUser)
		}
		if err != nil {
			panic(fmt.Sprintf("%s", err.Error()))
		}
	}
	// for user_detail
	for i := 0; i < _test_user_max_limit; i += 1 {
		newDetail := UserDetail{
			UserId:   uint64(_test_user_id_basic + i),
			Nickname: fmt.Sprintf("test_nickname_%d", i),
			Phone:    fmt.Sprintf("phone_%d", i),
			Email:    fmt.Sprintf("email_%d", i),
			Avatar:   fmt.Sprintf("avatar_%d", i),
		}
		var err error
		if i%2 == 0 {
			err = InsertWithName(ctx, _test_mysql_db, _test_sql_list[3], newDetail)
		} else {
			err = InsertWithPlace(ctx, _test_mysql_db, _test_sql_list[2], newDetail.UserId, newDetail.Nickname, newDetail.Phone, newDetail.Email, newDetail.Avatar)
		}
		if err != nil {
			panic(fmt.Sprintf("%s", err.Error()))
		}
	}
}

/*
	update 和 delete 直接下放到实际的程序执行时测试, 因为逻辑和insert一致, 几乎不可能出现预期外的结果
*/
// pass? -- yes just like insert happen
func Test_update_one(t *testing.T) {
}

// pass? -- yes just like insert happen
func Test_delete_one(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
	ctx, cancle := context.WithTimeout(context.Background(), _test_ctx_tomeout)
	defer cancle()
	_test_delall_basic_sql := `
	delete from user_basic
	`
	_test_delall_detail_sql := ` 
	delete from user_detail
	`
	var err error
	err = DeleteWithPlace(ctx, _test_mysql_db, _test_delall_detail_sql)
	err = DeleteWithPlace(ctx, _test_mysql_db, _test_delall_basic_sql)
	if err != nil {
		panic(err.Error())
	}
}

/*
实际上 select 也几乎不可能出现预期外的结果, 但是由于存在可能的model类型问题, 还是需要显式实现最小可执行测试进行部分结果显示测试.
1. one pass
2. more pass
3. page pass
*/
func Test_select_one(t *testing.T) {
	_test_mysql_db = InitMysql(_test_default_dsn)
	ctx, cancle := context.WithTimeout(context.Background(), _test_ctx_tomeout)
	defer cancle()
	// one test
	userBasic := UserBasic{}
	userDetail := UserDetail{}
	var err error
	userBasic, err = QueryWithPlace[UserBasic](ctx, _test_mysql_db, _test_sql_list[9], "zhangsan")
	userDetail, err = QueryWithName[UserDetail](ctx, _test_mysql_db, _test_sql_list[10], map[string]any{
		"user_id": 2,
	})
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	t.Logf("<=========================>\n")
	t.Logf("%+v\n", userBasic)
	t.Logf("%+v\n", userDetail)
	// more test
	userBasicList, err := QListWithPlace[UserBasic](ctx, _test_mysql_db, _test_sql_list[11], "test_user_%")
	userDetailList, err := QListWithName[UserDetail](ctx, _test_mysql_db, _test_sql_list[12], map[string]any{
		"min_user_id": 2,
	})
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	t.Logf("<=========================>\n")
	t.Logf("%+v\n", userBasicList)
	t.Logf("%+v\n", userDetailList)
	// page
	_test_current_page := 2
	_test_pagesize := 3
	page, err := QPageWithPlace[UserDetail](ctx, _test_mysql_db, _test_sql_list[13], (_test_current_page-1)*_test_pagesize, _test_pagesize)
	page.CurrentPage = _test_current_page
	page.PageSize = _test_pagesize
	page.Total = len(page.Items)
	t.Logf("<=========================>\n")
	t.Logf("%+v\n", page)
}
