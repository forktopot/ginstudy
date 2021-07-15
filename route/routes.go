package route

import (
	"ginstudy/controller"
	"ginstudy/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/index", controller.Index)
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", middleware.Session(), controller.Login)
	r.POST("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/auth/upload", controller.Upload)
	r.GET("/api/auth/download", controller.HandleDownloadFile)

	r.GET("/captcha", middleware.Session(), controller.Captcha)
	// r.GET("/captcha/verify/:value", middleware.Session(), controller.CaptchaVerify)

	return r
}
