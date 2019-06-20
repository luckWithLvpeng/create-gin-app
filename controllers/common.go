package controllers

//Response 返回结果
type Response struct {
	Msg  string
	Code int
	Data interface{}
}

const (
	//Success 成功
	Success = 200
	//Error 失败
	Error = 500
	// InvalidParams  无效的参数
	InvalidParams = 400
)

var mapMsg = map[int]string{
	Success:       "ok :",
	Error:         "failed :",
	InvalidParams: "请求参数错误 :",
}

// GetMsg   返回错误信息
func GetMsg(code int) string {
	msg, ok := mapMsg[code]
	if ok {
		return msg
	}
	return mapMsg[Error]
}
