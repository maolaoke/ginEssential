package common

import (
	"fmt"
	"ginEssential/model"
	"ginEssential/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var DB *gorm.DB = InitDB()

func InitDB() *gorm.DB {
	util.InitConfig()
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
	username,
	password,
	host,
	port,
	database,
	charset)
	fmt.Println(args)
	
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}) //自动创建表
	return db
}

func GetDB() *gorm.DB{
	return DB
}