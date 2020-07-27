package controllers

import (
	"github.com/astaxie/beego"
	"school-web/service"
)

type TokenController struct {
	beego.Controller
}

// GetToken ...
// @Title GetToken
// @Description GetToken
// @Param	grant_type		formData	string	true	"grant_type: password,refresh_token ..."
// @Param	username		formData	string	true	"username"
// @Param	password		formData	string	true	"password"
// @Param	scope   		formData	string	true	"read"
// @Param	client_id   	formData	string	true	"12345"
// @Param	client_secret   formData	string	true	"123456"
// @Success 200 {object} common.Result
// @Failure 403
// @router /token [post]
func (t *TokenController) GetToken() {
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
