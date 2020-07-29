package myError

import (
	"fmt"
	"runtime"
)

type Error struct {
	msg   string
	where string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.where, e.msg)
}

func New(format string, args ...interface{}) *Error {
	where := getWhere()
	format = getFormat(format, args)
	return &Error{
		msg:   format,
		where: where,
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
	var where string
	format = getFormat(format, args)
	switch t := err.(type) {
	case *Error:
		// 继承where
		where = t.where
		// 拼接上之前的错误
		format = t.msg + " -> " + format
	default:
		where = getWhere()
		format = format + " -> " + err.Error()
	}
	return &Error{
		msg:   format,
		where: where,
	}
}
