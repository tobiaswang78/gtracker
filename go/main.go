package main

import (
	bug "itflow/app"
	"itflow/app/bugconfig"
	"itflow/db"

	"github.com/hyahm/goconfig"
	"github.com/hyahm/golog"
)

func main() {

	// 初始化配置文件
	goconfig.InitConf("bug.ini", goconfig.INI)
	//初始化mysql
	db.InitMysql()
	// 初始化redis
	db.InitRedis()
	// 初始化缓存（后面会使用redis）
	bugconfig.LoadConfig()
	// 初始化日志
	golog.InitLogger(goconfig.ReadString("log.path", ""),
		goconfig.ReadInt64("log.size", 0),
		goconfig.ReadBool("log.everyday", false))

	////
	bug.RunHttp()
}
