package main

import (
	"photo_service/controllers"

	"photo_service/db"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDatabase() //初始化链接要在注册路由器之前
	ginServer := gin.Default()

	// ginServer.GET("/image", controllers.ImageHandler)
	// ginServer.GET("/readFile", controllers.ReadFile)

	//localhost:8080/encrypt?filename=images2.jfif
	ginServer.GET("/encrypt", controllers.EncryptImage)

	// http://localhost:8080/decrypt?id=3
	// ginServer.GET("/decrypt", controllers.DecryptImage)
	ginServer.GET("/decrypt", controllers.DecryptAndWatermark)

	// 将本地 "./storage" 文件夹映射为 "/static" 路由前缀

	ginServer.Run(":8080")

}
