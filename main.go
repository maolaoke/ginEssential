package main

import (
	"ginEssential/common"
	"ginEssential/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)


type User struct{
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size(255);not null"`
}


func main(){
	util.InitConfig()
	db := common.GetDB()
	
	defer db.Close()
	var r *gin.Engine = gin.Default()
	r = CollectRoute(r)

	port := viper.GetString("server.port")
	if port != ""{
		panic(r.Run(":"+port))
	}
	panic(r.Run())

}


