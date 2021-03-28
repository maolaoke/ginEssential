package common

import (
	"ginEssential/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)


var jwtKey = []byte("afghjruna")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func GetToken(user model.User)(string ,error){
	expirationTime := time.Now().Add(7*24*time.Hour)  //设置过期时间：7天
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "oceanlearn.tech",
			Subject: "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	if tokenString, err:= token.SignedString(jwtKey); err != nil{
		return "", err
	}else{
		return tokenString,nil
	}
}

//从tokenString中解析出相关信息
func ParseToken(tokenString string)(*jwt.Token, *Claims, error){
	claims := &Claims{}

	token,err := jwt.ParseWithClaims(tokenString, claims,
		 func(token *jwt.Token)(i interface{}, err error){
			return jwtKey, nil
		})
	return token, claims, err
}