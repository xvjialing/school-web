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
