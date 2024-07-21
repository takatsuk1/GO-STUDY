package main

import (
	"gin-mall/config"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/cache"
	"gin-mall/repository/db/dao"
	"gin-mall/routes"
)

// http:/127.0.0.1:3030/
func main() {
	loading()
	r := routes.NewRouter()
	_ = r.Run(":3000")
}

func loading() {
	config.InitConfig()
	dao.InitMySQL()
	log.InitLog()
	cache.InitCache()
}
