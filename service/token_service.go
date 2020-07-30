package service

import (
	"encoding/json"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"log"
	"net/http"
	"school-web/common"
	models2 "school-web/models"
	"strconv"
	"time"
)

var Srv *server.Server

func InitOauth2Service(adminUserName, adminUserPassword, adminUserEmail string) {

	initAdminUser(adminUserName, adminUserPassword, adminUserEmail)

	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	//tokenStore := mysql.NewDefaultStore(
	//	mysql.NewConfig("root:Xjl1994920!@tcp(rm-wz9t41ublf3gmrv3e5o.mysql.rds.aliyuncs.com:3306)/go-admin?charset=utf8"),
	//)
	//manager.MapTokenStorage(tokenStore)

	clientStore := store.NewClientStore()
	clientID := "12345"
	clientSecret := "123456"
	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)
	Srv = server.NewDefaultServer(manager)
	Srv.SetAllowGetAccessRequest(true)

	Srv.SetClientInfoHandler(server.ClientFormHandler)

	Srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {

		user, err := models2.CheckPassword(username, password)
		if user == nil {
			return "", err
		}
		return strconv.Itoa(user.Id), err
	})

	Srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	Srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
}

// HandleTokenRequest token request handling
func HandleTokenRequest(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	gt, tgr, err := s.ValidationTokenRequest(r)
	if err != nil {
		return tokenError(s, w, err)
	}

	ti, err := s.GetAccessToken(ctx, gt, tgr)
	if err != nil {
		return tokenError(s, w, err)
	}

	return token(s, w, s.GetTokenData(ti), nil)
}

func ValidationBearerToken(s *server.Server, w http.ResponseWriter, r *http.Request) error {

	token_info, err := s.ValidationBearerToken(r)

	if err != nil {
		return tokenError(s, w, err)
	}

	data := map[string]interface{}{
		"expires_in": int64(token_info.GetAccessCreateAt().Add(token_info.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token_info.GetClientID(),
		"user_id":    token_info.GetUserID(),
	}

	return token(s, w, data, nil)

}

func tokenError(s *server.Server, w http.ResponseWriter, err error) error {
	data, statusCode, header := s.GetErrorData(err)
	return token(s, w, data, header, statusCode)
}

func token(s *server.Server, w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	for key := range header {
		w.Header().Set(key, header.Get(key))
	}

	status := http.StatusOK
	w.WriteHeader(status)

	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	if status == http.StatusOK {
		return json.NewEncoder(w).Encode(common.Succes(data))
	} else {
		return json.NewEncoder(w).Encode(common.Unauthorized(data))
	}

}

func initAdminUser(adminUserName, adminUserPassword, adminUserEmail string) {

	exist := models2.ExistUserByUserName(adminUserName)

	if exist {
		return
	}

	user := models2.User{
		Username: adminUserName,
		Password: adminUserPassword,
		Email:    adminUserEmail,
	}

	id, err := models2.AddUser(&user)
	if err == nil {
		user.Id = int(id)
		log.Println("初始化用户id为", user)
	} else {
		log.Println("初始化用户失败", err)
	}
}

func CheckAccessToken(r *http.Request, access_token string) (tokenInfo oauth2.TokenInfo, err error) {

	tokenInfo, err = Srv.Manager.LoadAccessToken(r.Context(), access_token)

	return tokenInfo, err
}
