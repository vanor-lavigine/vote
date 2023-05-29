package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// Db 数据库句柄
	Db  *sql.DB
	err error
)

func InitDb() bool {
	Db, err = sql.Open("mysql", "debian-sys-maint:XMQWnyGB6Or12Oxk@tcp(localhost:3306)/vote")
	if err := Db.Ping(); err != nil {
		return false
	}
	return true
}
