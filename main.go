package main

import (
	"ginstudy/common"
	"ginstudy/route"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	InitConfig() //读取配置
	db := common.InitDB()
	defer db.Close()
	// 1.创建路由
	r := gin.Default()
	r.LoadHTMLGlob("./*.html")

	r = route.CollectRoute(r)

	// r.Use(middleware.AuthMiddleware())
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	port := viper.GetString("server.port")
	if port == "" {
		panic(r.Run(":8000"))
	}
	panic(r.Run(":" + port))
}

func InitConfig() {
	workDir, _ := os.Getwd()                 //获取当前根目录
	viper.SetConfigName("application")       //设置配置文件名
	viper.SetConfigType("yml")               //设置配置文件类型
	viper.AddConfigPath(workDir + "/config") //设置配置文件路径
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
