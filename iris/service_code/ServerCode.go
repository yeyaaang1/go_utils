package service_code

const (
	Success      = 1
	UnknownError = -1

	ParamsVerifyError = 8101
	MethodNotAllowed  = 8201
)

var ServerMsgMap = map[int]string{
	Success:      "成功",
	UnknownError: "未知错误",

	ParamsVerifyError: "参数验证失败",
	MethodNotAllowed:  "不被允许的方法",
}

func GetServerMsg(code int) string {
	msg, ok := ServerMsgMap[code]
	if ok {
		return msg
	}
	return "unknown service code"
}
