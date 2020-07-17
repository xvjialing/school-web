package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
	"school-web/service"
	"time"
)

type TokenController struct {
	beego.Controller
}

// GetToken ...
// @Title GetToken
// @Description GetToken
// @Param	grant_type	query	string	false	"grant_type. e.g. col1:v1,col2:v2 ..."
// @Param	username	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} common.Result
// @Failure 403
// @router /token [post]
func (t *TokenController) GetToken() {
	service.HandleTokenRequest(service.Srv, t.Ctx.ResponseWriter, t.Ctx.Request)
}

// @router /info [post]
func (t *TokenController) GetTokenInfo() {
	token, err := service.Srv.ValidationBearerToken(t.Ctx.Request)
	if err != nil {
		http.Error(t.Ctx.ResponseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(t.Ctx.ResponseWriter)
	e.SetIndent("", "  ")
	e.Encode(data)
}
