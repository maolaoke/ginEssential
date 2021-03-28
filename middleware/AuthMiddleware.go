package middleware

import (
	"ginEssential/common"
	"ginEssential/model"
	"net/http"
	"strings"
	"ginEssential/dto"
	"github.com/gin-gonic/gin"
)

//验证解析token
func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//验证token格式,若token为空或不是以Bearer开头，则token格式不对
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer"){
			ctx.JSON(http.StatusUnauthorized, gin.H{"code":401, "msg":"权限不足"})
			ctx.Abort()  //将此次请求抛弃
			return
		}

		tokenString = tokenString[7:]  //token的前面是“bearer”，有效部分从第7位开始

		//从tokenString中解析信息
		token, claims, err:=common.ParseToken(tokenString)
		if err != nil || !token.Valid{
			ctx.JSON(http.StatusUnauthorized, gin.H{"code":401, "msg":"权限不足"})
			ctx.Abort()
			return
		}

		// 查询tokenString中的user信息是否存在
		userId := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userId)

		if user.ID == 0{
			ctx.JSON(http.StatusUnauthorized, gin.H{"code":401, "msg":"权限不足"})
			ctx.Abort()
			return
		}

		//若存在该用户则将用户信息写入上下文
		userDto := dto.ToUserDto(&user)
		ctx.Set("user", userDto)
		ctx.Next()
	}
}