package main

import (
	"goIM/initialize"
	"goIM/router"
)

func main() {
	// 初始化日志
	initialize.InitLogger()
	// 初始化数据库
	initialize.InitDB()

	ginEngine := router.InitRouter()
	ginEngine.Run(":8000")
}
