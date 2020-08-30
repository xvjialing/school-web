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
	ctx.Output.Header("Access-Control-Allow-Headers", "X-Custom-Header,accept,Content-Type,Access-Token,access_token")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Origin", origin)
	if ctx.Input.Method() == http.MethodOptions {
		// options请求，返回200
		ctx.Output.SetStatus(http.StatusOK)
		_ = ctx.Output.Body(success)
	}
}

func main() {

	beegoAppConfig := beego.AppConfig

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("sqlconn"))
	//orm.Debug = true

	beego.InsertFilter("/*", beego.BeforeRouter, corsFunc)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	adminUserName := beegoAppConfig.String("AdminUserName")
	adminUserPassword := beegoAppConfig.String("AdminUserPassword")
	adminUserEmail := beegoAppConfig.String("AdminUserEmail")
	RedisAddr := beegoAppConfig.String("RedisAddr")
	RedisDB, _ := beegoAppConfig.Int("RedisDB")
	RedisPassword := beegoAppConfig.String("RedisPassword")

	service.InitOauth2Service(adminUserName, adminUserPassword, adminUserEmail, RedisAddr, RedisPassword, RedisDB)

	beego.Run()
}
