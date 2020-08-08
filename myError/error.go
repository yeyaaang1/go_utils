package myError

import (
	"fmt"
	"runtime"
)

type Error struct {
	msg   string
	where []string
}

func (e *Error) Error() string {
	return fmt.Sprintf("{\"where\":\"%v\", \"msg\": \"%s\"}",
		e.where, e.msg)
}

func (e *Error) GetMsg() string {
	return e.msg
}

func (e *Error) GetWhere() string {
	return e.msg
}

func New(format string, args ...interface{}) *Error {
	where := getWhere()
	format = getFormat(format, args)
	return &Error{
		msg:   format,
		where: []string{where},
	}
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
		}
	}
	var whereSlice []string
	switch t := err.(type) {
	case *Error:
		// 继承where
		whereSlice = append(t.where, where)
		// 拼接上之前的错误
		format = t.msg + " -> " + format
	default:
		whereSlice = []string{where}
		format = format + " -> " + err.Error()
	}
	return &Error{
		msg:   format,
		where: whereSlice,
	}
}
