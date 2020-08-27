package controllers

import (
	"github.com/astaxie/beego"
	"school-web/service"
)

//用户认证相关操作
type TokenController struct {
	beego.Controller
}

// 登录获取token ...
// @Title GetToken
// @Description GetToken
// @Param	grant_type		formData	string	true	"grant_type: password,refresh_token ..."
// @Param	username		formData	string	true	"username"
// @Param	password		formData	string	true	"password"
// @Success 200 {object} common.Result
// @Failure 403
// @router /token [post]
func (t *TokenController) GetToken() {
	t.Ctx.Request.Form.Add("scope", "read")
	t.Ctx.Request.Form.Add("client_id", "12345")
	t.Ctx.Request.Form.Add("client_secret", "123456")
	service.HandleTokenRequest(service.Srv, t.Ctx.ResponseWriter, t.Ctx.Request)
}

// GetTokenInfo ...
// @Title GetTokenInfo
// @Description GetToken
// @Param	access_token		formData	string	true	"access_token"
// @Success 200 {object} common.Result
// @Failure 403
// @router /info [post]
func (t *TokenController) GetTokenInfo() {
	service.ValidationBearerToken(service.Srv, t.Ctx.ResponseWriter, t.Ctx.Request)
}
