package restful

import (
	"encoding/json"
	"gitee.com/super_step/go_utils/iris/service_code"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
)

type RestData struct {
	Msg  string
	Code int
	Data interface{}
}

func OK(data interface{}, log ...bool) mvc.Result {
	return Result(http.StatusOK, RestData{
		Msg:  service_code.GetServerMsg(service_code.Success),
		Code: service_code.Success,
		Data: data,
	}, nil, log...)
}

func ParamsError(err error, log ...bool) mvc.Result {
	return Result(http.StatusBadRequest, RestData{
		Msg:  service_code.GetServerMsg(service_code.ParamsVerifyError),
		Code: service_code.ParamsVerifyError,
	}, err, log...)
}

func MethodError(log ...bool) mvc.Result {
	return Result(http.StatusMethodNotAllowed, RestData{
		Msg:  service_code.GetServerMsg(service_code.MethodNotAllowed),
		Code: service_code.MethodNotAllowed,
	}, nil, log...)
}

func UnAuth(code int, err error, log ...bool) mvc.Result {
	return Result(http.StatusUnauthorized, RestData{
		Msg:  service_code.GetServerMsg(code),
		Code: code,
		Data: nil,
	}, err, log...)
}

func BadRequest(code int, err error, log ...bool) mvc.Result {
	return Result(http.StatusBadRequest, RestData{
		Msg:  service_code.GetServerMsg(code),
		Code: code,
		Data: nil,
	}, err, log...)
}

func Result(httpCode int, restData RestData, err error, log ...bool) mvc.Result {
	result := iris.Map{
		"msg":  restData.Msg,
		"code": restData.Code,
	}
	if restData.Data != nil {
		result["data"] = restData.Data
	}
	if err != nil {
		golog.Default.Error(err.Error())
		result["error"] = err.Error()
	}
	var logout bool
	if len(log) > 0 {
		logout = log[0]
	}
	if httpCode >= 300 || logout {
		var logFunc func(format string, args ...interface{})
		if httpCode >= 300 {
			logFunc = golog.Default.Warnf
		} else {
			logFunc = golog.Default.Infof
		}
		resJson, err := json.Marshal(result)
		if err != nil {
			logFunc("response: %+v", result)
		}
		logFunc("response: %s", resJson)
	}
	return mvc.Response{
		Code:   httpCode,
		Object: result,
	}
}
