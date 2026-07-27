package main

import (
	"Ielts-vocab/internal/database"
	"Ielts-vocab/internal/server"
	"log"

	"github.com/gin-gonic/gin"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// 初始化数据库
	db, err := database.ConnectToPostgres()
	if err != nil {
		log.Fatal(err)
	}

	// 创建 gin
	r := gin.Default()

	// 注册路由
	server.SetupRouter(r, db)

	log.Fatal(r.Run(":8080"))
}
