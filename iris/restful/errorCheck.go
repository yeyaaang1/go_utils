package restful

import (
	"gitee.com/super_step/go_utils/myError"
)

func ErrorFormat(err error) string {
	switch t := err.(type) {
	case *myError.Error:
		return t.GetMsg()
	default:
		return t.Error()
	}
}
