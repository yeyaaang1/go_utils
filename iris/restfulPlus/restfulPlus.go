package restfulPlus

import (
	"gitee.com/super_step/go_utils/iris/service_code"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"net/http"
)

type RestData struct {
	Msg  string
	Code int
	Data interface{}
}

func OK(ctx iris.Context, data interface{}) mvc.Result {
	return JsonResult(ctx, http.StatusOK, RestData{
		Msg:  service_code.GetServerMsg(service_code.Success),
		Code: service_code.Success,
		Data: data,
	}, nil)
}

func ParamsError(ctx iris.Context, err error) mvc.Result {
	return JsonResult(ctx, http.StatusBadRequest, RestData{
		Msg:  service_code.GetServerMsg(service_code.ParamsVerifyError),
		Code: service_code.ParamsVerifyError,
	}, err)
}

func MethodError(ctx iris.Context) mvc.Result {
	return JsonResult(ctx, http.StatusMethodNotAllowed, RestData{
		Msg:  service_code.GetServerMsg(service_code.MethodNotAllowed),
		Code: service_code.MethodNotAllowed,
	}, nil)
}

func UnAuth(ctx iris.Context, code int, err error) mvc.Result {
	return JsonResult(ctx, http.StatusUnauthorized, RestData{
		Msg:  service_code.GetServerMsg(code),
		Code: code,
		Data: nil,
	}, err)
}

func BadRequest(ctx iris.Context, code int, err error) mvc.Result {
	return JsonResult(ctx, http.StatusBadRequest, RestData{
		Msg:  service_code.GetServerMsg(code),
		Code: code,
		Data: nil,
	}, err)
}

func JsonResult(ctx iris.Context, httpCode int, restData RestData, err error) mvc.Result {
	result := iris.Map{
		"msg":  restData.Msg,
		"code": restData.Code,
	}
	if restData.Data != nil {
		result["data"] = restData.Data
	}
	if err != nil {
		result["error"] = err.Error()
	}
	// todo 在中间件中实现log功能
	return mvc.Response{
		Code:   httpCode,
		Object: result,
	}
}
