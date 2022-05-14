package db

import (
	"github.com/jinzhu/gorm"
	// コメント書かないとLintエラーになる
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Init() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@/go?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect db")
	}
	db.LogMode(false)
	return db
}
