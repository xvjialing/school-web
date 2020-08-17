package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	_ "school-web/routers"
	"school-web/service"
)

func init() {
	log.SetPrefix("[INFO]")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

var success = []byte("SUPPORT OPTIONS")

var corsFunc = func(ctx *context.Context) {
	origin := ctx.Input.Header("Origin")
	ctx.Output.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
	ctx.Output.Header("Access-Control-Max-Age", "3600")
	ctx.Output.Header("Access-Control-Allow-Headers", "X-Custom-Header,accept,Content-Type,Access-Token")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Origin", origin)
	if ctx.Input.Method() == http.MethodOptions {
		// options请求，返回200
		ctx.Output.SetStatus(http.StatusOK)
		_ = ctx.Output.Body(success)
	}
}

func main() {

	//beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	//	//允许访问所有源
	//	AllowAllOrigins: true,
	//	//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
	//	//其中Options跨域复杂请求预检
	//	AllowMethods: []string{"*"},
	//	//指的是允许的Header的种类
	//	AllowHeaders: []string{"*"},
	//	//公开的HTTP标头列表
	//	ExposeHeaders: []string{"Accept", "Origin", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	//	//如果设置，则允许共享身份验证凭据，例如cookie
	//	AllowCredentials: true,
	//}))

	beego.InsertFilter("/*", beego.BeforeRouter, corsFunc)

	beegoAppConfig := beego.AppConfig

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("sqlconn"))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	adminUserName := beegoAppConfig.String("AdminUserName")
	adminUserPassword := beegoAppConfig.String("AdminUserPassword")
	adminUserEmail := beegoAppConfig.String("AdminUserEmail")

	service.InitOauth2Service(adminUserName, adminUserPassword, adminUserEmail)

	beego.Run()
}
