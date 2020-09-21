package filter

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
	"log"
)

func LogRequest(ctx *context.Context) {

	url, _ := json.Marshal(ctx.Input.Data()["RouterPattern"])
	params, _ := json.Marshal(ctx.Request.Form)
	outputBytes, _ := json.Marshal(ctx.Input.Data()["json"])
	outputStr := "requestUrl:" + string(url) + ", requestParams:" + string(params) + ", response:" + string(outputBytes)
	log.Println(outputStr)
}
