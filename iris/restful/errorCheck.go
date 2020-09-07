package restful

import (
	"gitee.com/super_step/go_utils/iris/service_code"
	"gitee.com/super_step/go_utils/myError"
)

func ErrorFormat(err error) (code int, msg string) {
	switch t := err.(type) {
	case *myError.Error:
		return t.GetCode(), t.GetMsg()
	default:
		return service_code.UnknownError, t.Error()
	}
}
