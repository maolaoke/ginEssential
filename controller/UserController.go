package controller

import (
	"ginEssential/common"
	"ginEssential/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)


func Register(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证

	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须11位"})
		return
	}
	if len(password) < 8 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码必须大于8位"})
		return
	}
	log.Println(name, telephone, password)

	//判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号已注册"})
		return
	}

	//创建用户
	hashdPassword,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)  //密码hash化
	if err != nil{
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashdPassword),
	}
	if err := db.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": err})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg": "注册成功",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool{
	var user model.User
	db.Where("telephone=?",telephone).First(&user)
	return user.ID != 0
}

func Login(ctx *gin.Context){
	db := common.GetDB()
	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据校验
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须11位"})
		return
	}
	if len(password) < 8 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码必须大于8位"})
		return
	}

	//判断手机号是否存在
	var user model.User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID == 0{
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号不存在"})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发送token
	token, err := common.GetToken(user)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"code":500, "msg":"系统token获取异常"})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"token":token,
		},
		"msg": "登录成功",
	})

}


func Info(ctx *gin.Context){
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code":200, "data":gin.H{"user":user}})
}