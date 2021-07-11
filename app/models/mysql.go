package models

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

var ChatDB gdb.DB

func InitDB() {
	// 參考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
	// dsn := viper.GetString(`mysql.dsn`)
	// ChatDB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	ChatDB = g.DB()

}
