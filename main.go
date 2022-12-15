package main

import (
	"log"
	"mycloud/db/mysql"
	"mycloud/lib"
	"mycloud/router"
)

func main() {
	serverConfig := lib.LoadServerConfig()
	mysql.InitDB(serverConfig)
	defer mysql.DB.Close()

	r := router.SetupRoute()

	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")

	if err := r.Run(":8010"); err != nil {
		log.Fatal("服务器启动失败...")
	}
}
