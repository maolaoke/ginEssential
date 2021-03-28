package common

import (
	"ginEssential/model"
	
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB = initDB()

func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/ginessential?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}) //自动创建表
	return db
}

func GetDB() *gorm.DB{
	return DB
}