## go+gin+grom+JWT+viper实现登录注册功能

### 1 如何gin创建server监听请求

> 安装gin

```
go get github.com/gin-gonic/gin
```

> 开启一个server

```go
var r *gin.Engine = gin.Default()		// Engine is the framework's instance
r.POST("/register", controller.Register)  //监听post请求。默认是8080端口
panic(r.Run())   //启动程序，此时就可以访问localhost:8080/register
```

其中的controller.Register函数如下：

```go
func Register(ctx *gin.Context) {
	//获取请求body参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须11位"})
		return
	}

	//判断手机号是否存在
	//创建用户
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg": "注册成功",
	})
}
```





### 2 gorm操作数据库

> 安装gorm以及mysql驱动

```
go get github.com/jinzhu/gorm
go get github.com/jinzhu/gorm/dialects/mysql
```

> 配置数据库信息

```go
import(
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"  //mysql驱动
)

//初始化，返回一个*gorm.DB类型的db，以后操作数据库都用db去操作
func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{}) //自动创建表，默认会添加id/created_at/updated_at/deleted_at等字段
	return db
}

//User实体类
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size(255);not null"`
}
```

> CRUD

1. 添加记录

   ```go
   db.Create(&user)  //方式一
   
   db.NewRecord(user)  //方式二
   ```

2. 删除记录

   ```go
   db.Delete(&user)  //delete from users where id =10
   
   db.Where("email LIKE ?", "%@qq.com").Delete(User{})   //delete from users where email like "%@qq.com"
   
   // 有delete_at字段时，默认软删除,即不会真正删除数据，而是设置删除时间
   ```

3. 修改记录

   ```go
   db.Save(&user)  //修改所有字段
   
   db.Model(&user).Update("name", "maolaoke")  //修改单个字段
   ```

4. 查询记录

   ```go
   db.First(&user) 	//查询第一条记录，按主键排序
   db.Last(&user)
   
   db.Find(&user) 		//查询所有记录
   
   db.Where("name=?","maolaoke").First(&user) 	//SELECT * FROM users WHERE name = 'maolaoke' limit 1
   ```

   



### 3 JWT的token认证

1. 安装jwt工具包

   ```bash
   go get github.com/dgrijalva/jwt-go
   ```

2. 创建一个结构体存储保存在token中的数据

   ```go
   type Claims struct {
   	UserId uint
   	jwt.StandardClaims   //标准的token应有的信息（如创建时间，过期时间等）
   }
   ```

3. 创建获取token字符串的函数

   ```go
   func GetToken(user model.User)(string ,error){
   	expirationTime := time.Now().Add(7*24*time.Hour)  //设置过期时间：7天
   	claims := &Claims{   //创建Claims实例
   		UserId: user.ID,
   		StandardClaims: jwt.StandardClaims{
   			ExpiresAt: expirationTime.Unix(),
   			IssuedAt: time.Now().Unix(),
   		},
   	}
   
   	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)  //根据claims创建token
   	if tokenString, err:= token.SignedString(jwtKey); err != nil{
   		return "", err
   	}else{
   		return tokenString,nil  //返回token字符串
   	}
   }
   ```

4. 在登录成功后返回token字符串

   ```go
   token, err := common.GetToken(user)
   if err != nil {
       response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统token获取异常")
       return
   }
   response.Success(ctx, gin.H{"token": token}, "登录成功")
   ```

5. 在网站中我们经常需要根据token免登录获取用户信息。

   1. 创建中间验证解析token函数

      ```go
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
      		userDto := dto.ToUserDto(&user)   //UserDto是为了防止所有user信息返回到前端
      		ctx.Set("user", userDto)
      		ctx.Next()
      	}
      }
      ```

      ```go
      //从tokenString中解析出相关信息
      func ParseToken(tokenString string)(*jwt.Token, *Claims, error){
      	claims := &Claims{}
      
      	token,err := jwt.ParseWithClaims(tokenString, claims,
      		 func(token *jwt.Token)(i interface{}, err error){
      			return jwtKey, nil
      		})
      	return token, claims, err
      }
      ```

      

      ```go
      type UserDto struct {
      	Name      string `json:"name"`    //由于返回前端的变量一般都是小写开头，这里规范一下
      	Telephone string `json:"telephone"`
      }
      
      func ToUserDto(user *model.User)UserDto{
      	return UserDto{
      		Name: user.Name,
      		Telephone: user.Telephone,
      	}
      }
      ```

      

   2. 在路由中添加中间拦截器

      ```go
      r.GET("api/auth/info",middleware.AuthMiddleware(), controller.Info)
      ```

   3. 创建controller

      ```go
      func Info(ctx *gin.Context) {
      	user, _ := ctx.Get("user")  //若有正确的token，则1中的函数会将user信息写入ctx，这里直接获取就行
      	response.Success(ctx, gin.H{"user": user}, "")
      }
      ```

      



### 4 统一封装格式

在返回到前端的json中，一般都是统一格式。

```go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//统一返回格式
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 200, data, msg)
}

func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 400, data, msg)
}
```



### 5 配置文件

使用配置文件来管理所有的配置

1. 下载viper包

   ```go
   go get github.com/spf13/viper
   ```

2. 创建配置文件`config/application.yml`

   ```yml
   server:
     port: 8080
   
   datasource:
     driverName: mysql
     host: 127.0.0.1
     port: 3306
     database: ginessential
     username: root
     password: root
     charset: utf8
   ```

   

3. 在viper中设置配置文件路径即格式

   ```go
   func InitConfig() {
   	workDir, _ := os.Getwd() //获取当前目录
   	viper.SetConfigName("application")
   	viper.SetConfigType("yml")
   	viper.AddConfigPath(workDir + "/config")
   	err := viper.ReadInConfig()
   	if err != nil {
   		panic(err)
   	}
   }
   ```

4. 使用viper获取配置信息

   ```go
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
   ```

   





