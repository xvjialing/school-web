package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	"log"
	_ "school-web/routers"
	"school-web/service"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	log.SetPrefix("[INFO]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//允许访问所有源
		AllowAllOrigins: true,
		//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		//其中Options跨域复杂请求预检
		AllowMethods: []string{"*"},
		//指的是允许的Header的种类
		AllowHeaders: []string{"*"},
		//公开的HTTP标头列表
		ExposeHeaders: []string{"Content-Length"},
		//如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("sqlconn"))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	service.InitOauth2Service()

	beego.Run()
}
