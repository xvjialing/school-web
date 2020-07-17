package common

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Succes(data interface{}) (result Result) {
	result.Code = 200
	result.Msg = "ok"
	result.Data = data

	return result
}

func Unauthorized(data interface{}) (result Result) {
	result.Code = 401
	result.Msg = "unauthorized"
	result.Data = data
	return result
}

func Failed(code int, msg string) (result Result) {

	result.Code = code
	result.Msg = msg

	return result
}
