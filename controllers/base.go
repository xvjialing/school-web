package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"log"
	"school-web/common"
	"school-web/models"
	"school-web/service"
	"strconv"
)

//const PrometheusUrl string = "/v1/prometheus"

type BaseController struct {
	beego.Controller
	User       *models.User
	RequestUri string
}

func (c *BaseController) Prepare() {
	requestUri := c.Ctx.Request.RequestURI

	log.Println("requestUri : " + requestUri)

	//if PrometheusUrl == requestUri {
	//	return
	//}

	access_token := c.Ctx.Input.Query("access_token")

	if len(access_token) == 0 {

		c.Data["json"] = common.Failed(401, "access_token can not be empty")
		c.ServeJSON()
		c.StopRun()
		return
	}
	var err error

	tokenInfo, err := service.CheckAccessToken(c.Ctx.Request, access_token)
	if err != nil {
		c.Data["json"] = common.Failed(401, "access_token is error")
		c.ServeJSON()
		c.StopRun()
		return
	}
	userId, _ := strconv.Atoi(tokenInfo.GetUserID())
	c.User, err = models.GetUserById(userId)

	if err != nil {
		c.Data["json"] = common.Failed(401, "Authorization failed")
		c.ServeJSON()
		c.StopRun()
		return
	}

	bytes, _ := json.Marshal(c.User)
	log.Println("user:", string(bytes))
}
