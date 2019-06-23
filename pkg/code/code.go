package code

const (
	//Success 成功
	Success = 200
	//Error 失败
	Error = 500
	//ErrorRecordNotFound 未找到值
	ErrorRecordNotFound = 404
	// InvalidParams  无效的参数
	InvalidParams = 400

	// AuthFail 权限验证失败
	AuthFail = 1000
	// AuthNeedHeaderAuthorization 请求headers 缺少 Authorization 信息
	AuthNeedHeaderAuthorization = 1001
	// AuthParseToken 解析token失败
	AuthParseToken = 1002
	// AuthTokenTimeout  token 过期
	AuthTokenTimeout = 1003
	// AuthInvalidUsernamePasssword  无效的用户名密码
	AuthInvalidUsernamePasssword = 1004
	//AuthInvalid  token 失效
	AuthInvalid = 1005
)

var mapMsg = map[int]string{
	Success:                      "ok",
	Error:                        "failed: ",
	InvalidParams:                "请求参数错误: ",
	ErrorRecordNotFound:          "未找到记录: ",
	AuthFail:                     "权限验证失败",
	AuthNeedHeaderAuthorization:  "缺少Authorization权限认证信息",
	AuthParseToken:               "解析Authorization信息失败",
	AuthTokenTimeout:             "权限已经过期",
	AuthInvalidUsernamePasssword: "无效的用户名和密码",
	AuthInvalid:                  "权限已失效",
}

// GetMsg   返回错误信息
func GetMsg(code int) string {
	msg, ok := mapMsg[code]
	if ok {
		return msg
	}
	return mapMsg[Error]
}
