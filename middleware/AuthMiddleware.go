package middleware

import (
	"ginstudy/common"
	"ginstudy/model"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//认证

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取authorization header
		tokenString := c.GetHeader("Authorization")

		//验证token格式(不为空且以 Bearer 开头)
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() //抛弃这一次请求
			return
		}

		tokenString = tokenString[7:] //提取token有效字符，Bearer 之后的字符

		//验证token是否正确
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() //抛弃这一次请求
			return
		}

		//验证通过后获取 claim 中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//判断用户是否存在
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() //抛弃这一次请求
			return
		}

		//将user的信息写入上下文
		c.Set("user", user)

		c.Next()

	}
}

// func SessionConfig() sessions.Store {
// 	sessionMaxAge := 3600
// 	sessionSecret := "topgoer"
// 	var store sessions.Store
// 	store = cookie.NewStore([]byte(sessionSecret))
// 	store.Options(sessions.Options{
// 		MaxAge: sessionMaxAge, //seconds
// 		Path:   "/",
// 	})
// 	return store
// }

// 中间件，处理session
func Session() gin.HandlerFunc {
	sessionMaxAge := 3600
	sessionSecret := "topgoer"
	var store sessions.Store
	store = cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge: sessionMaxAge, //seconds
		Path:   "/",
	})
	return sessions.Sessions("topgoer", store)
}
