package main

import (
	"ginEssential/common"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)


type User struct{
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size(255);not null"`
}


func main(){
	db := common.GetDB()
	
	defer db.Close()
	var r *gin.Engine = gin.Default()
	r = CollectRoute(r)
	panic(r.Run())
}

