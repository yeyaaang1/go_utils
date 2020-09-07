package myError

import (
	"fmt"
	"gitee.com/super_step/go_utils/iris/service_code"
	"runtime"
)

type Error struct {
	msg   string
	where []string
	code  int
}

func (e *Error) Error() string {
	return fmt.Sprintf("{\"where\":\"%v\", \"msg\": \"%s\", \"code\": %d}",
		e.where, e.msg, e.code)
}

func (e *Error) GetMsg() string {
	return e.msg
}

func (e *Error) GetWhere() string {
	return e.msg
}

func (e *Error) GetCode() int {
	return e.code
}

func New(format string, args ...interface{}) *Error {
	where := getWhere()
	format = getFormat(format, args)
	return &Error{
		msg:   format,
		where: []string{where},
		code:  service_code.UnknownError,
	}
}

func NewWithCode(code int, format string, args ...interface{}) *Error {
	tmpErr := New(format, args...)
	tmpErr.code = code
	return tmpErr
}

func getWhere() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s(%d)", file, line)
}

func getFormat(format string, args []interface{}) string {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	return format
}

func Warp(err error, format string, args ...interface{}) *Error {
	where := getWhere()
	format = getFormat(format, args)
	if err == nil {
		return &Error{
			msg:   format,
			where: []string{where},
			code:  service_code.UnknownError,
		}
	}
	var (
		whereSlice []string
		code       int
	)
	switch t := err.(type) {
	case *Error:
		// 继承where
		whereSlice = append(t.where, where)
		// 拼接上之前的错误
		format = t.msg + " -> " + format
		code = t.code
	default:
		whereSlice = []string{where}
		format = format + " -> " + err.Error()
		code = service_code.UnknownError
	}
	return &Error{
		msg:   format,
		where: whereSlice,
		code:  code,
	}
}

func WarpWithCode(err error, code int, format string, args ...interface{}) *Error {
	tmpErr := Warp(err, format, args...)
	if tmpErr.code == service_code.UnknownError {
		tmpErr.code = code
	}
	return tmpErr
}
