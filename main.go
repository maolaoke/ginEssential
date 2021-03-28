package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


type User struct{
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size(255);not null"`
}


func main(){
	db := InitDB()
	defer db.Close()
	r:=gin.Default()
	r.POST("/register", func(ctx *gin.Context){
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		//数据验证
		fmt.Println(len(telephone))
		if(len(telephone) != 11){
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg":"手机号必须11位"})
			return
		}
		if(len(password) < 8){
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg":"密码必须大于8位"})
			return
		}
		log.Println(name, telephone, password)

		//判断手机号是否存在
		if isTelephoneExist(db, telephone){
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg":"手机号已注册"})
			return
		}

		//创建用户
		newUser :=User{
			Name:name,
			Telephone: telephone,
			Password: password,
		}
		if err:=db.Create(&newUser).Error;err != nil{
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg":err})
			return
		}

		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})

	panic(r.Run())
}

func isTelephoneExist(db *gorm.DB, telephone string) bool{
	var user User
	db.Where("telephone=?",telephone).First(&user)
	return user.ID != 0
}


func InitDB() *gorm.DB{
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/ginessential?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil{
		panic(err)
	}
	db.AutoMigrate(&User{})  //自动创建表
	return db
}