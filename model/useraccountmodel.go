package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userAccountFieldNames          = builder.RawFieldNames(&UserAccount{})
	userAccountRows                = strings.Join(userAccountFieldNames, ",")
	userAccountRowsExpectAutoSet   = strings.Join(stringx.Remove(userAccountFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	userAccountRowsWithPlaceHolder = strings.Join(stringx.Remove(userAccountFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheUserAccountIdPrefix     = "cache:userAccount:id:"
	cacheUserAccountUserIdPrefix = "cache:userAccount:userId:"
)

type (
	UserAccountModel interface {
		Insert(data *UserAccount) (sql.Result, error)
		FindOne(id int64) (*UserAccount, error)
		FindOneByUserId(userId int64) (*UserAccount, error)
		Update(data *UserAccount) error
		Delete(id int64) error
		TxAdjustBalance(tx *sql.Tx, userId int64, amount int64) (sql.Result, error)
	}

	defaultUserAccountModel struct {
		sqlc.CachedConn
		table string
	}

	UserAccount struct {
		Id         int64     `db:"id"`
		UserId     int64     `db:"user_id"`
		Balance    float64   `db:"balance"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewUserAccountModel(conn sqlx.SqlConn, c cache.CacheConf) UserAccountModel {
	return &defaultUserAccountModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_account`",
	}
}

func (m *defaultUserAccountModel) TxAdjustBalance(tx *sql.Tx, userId int64, amount int64) (sql.Result, error) {
	IdKey := fmt.Sprintf("%s%v", cacheUserAccountUserIdPrefix, userId)

	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set balance=balance+? where balance >= -? and user_id=?", m.table)
		fmt.Println(query, userId, amount)
		return tx.Exec(query, amount, amount, userId)
	}, IdKey)
}

func (m *defaultUserAccountModel) Insert(data *UserAccount) (sql.Result, error) {
	userAccountIdKey := fmt.Sprintf("%s%v", cacheUserAccountIdPrefix, data.Id)
	userAccountUserIdKey := fmt.Sprintf("%s%v", cacheUserAccountUserIdPrefix, data.UserId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, userAccountRowsExpectAutoSet)
		return conn.Exec(query, data.UserId, data.Balance)
	}, userAccountIdKey, userAccountUserIdKey)
	return ret, err
}

func (m *defaultUserAccountModel) FindOne(id int64) (*UserAccount, error) {
	userAccountIdKey := fmt.Sprintf("%s%v", cacheUserAccountIdPrefix, id)
	var resp UserAccount
	err := m.QueryRow(&resp, userAccountIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAccountRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAccountModel) FindOneByUserId(userId int64) (*UserAccount, error) {
	userAccountUserIdKey := fmt.Sprintf("%s%v", cacheUserAccountUserIdPrefix, userId)
	var resp UserAccount
	err := m.QueryRowIndex(&resp, userAccountUserIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userAccountRows, m.table)
		if err := conn.QueryRow(&resp, query, userId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAccountModel) Update(data *UserAccount) error {
	userAccountIdKey := fmt.Sprintf("%s%v", cacheUserAccountIdPrefix, data.Id)
	userAccountUserIdKey := fmt.Sprintf("%s%v", cacheUserAccountUserIdPrefix, data.UserId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAccountRowsWithPlaceHolder)
		return conn.Exec(query, data.UserId, data.Balance, data.Id)
	}, userAccountIdKey, userAccountUserIdKey)
	return err
}

func (m *defaultUserAccountModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	userAccountIdKey := fmt.Sprintf("%s%v", cacheUserAccountIdPrefix, id)
	userAccountUserIdKey := fmt.Sprintf("%s%v", cacheUserAccountUserIdPrefix, data.UserId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userAccountIdKey, userAccountUserIdKey)
	return err
}

func (m *defaultUserAccountModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserAccountIdPrefix, primary)
}

func (m *defaultUserAccountModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAccountRows, m.table)
	return conn.QueryRow(v, query, primary)
}
