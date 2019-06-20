package errorcode

import "net/http"

const (
	//ErrorExistTag 标签已经存在
	ErrorExistTag = 10000

	// ErrorAuth  auth 验证失败
	ErrorAuth = 20000
)

// MsgFlags  错误表
var MsgFlags = map[int]string{
	http.StatusInternalServerError: "fail",
	ErrorExistTag:                  "已存在",
	ErrorAuth:                      "权限验证失败",
}

// GetMsg 获取错误信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[http.StatusInternalServerError]
}
