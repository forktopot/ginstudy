package main

import (
	"ginstudy/common"
	"ginstudy/route"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := common.InitDB()
	defer db.Close()
	// 1.创建路由
	r := gin.Default()

	r = route.CollectRoute(r)
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	panic(r.Run(":8000"))
}
