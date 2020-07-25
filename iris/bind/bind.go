package bind

import (
	"encoding/json"
	"errors"
	"gitee.com/super_step/go_utils/myError"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"google.golang.org/protobuf/proto"
	"reflect"
	"runtime"
	"strings"

	zhCn "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator
var Validate *validator.Validate

func init() {
	trans, _ = ut.New(zhCn.New()).GetTranslator("zh")
	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	_ = zhTranslations.RegisterDefaultTranslations(Validate, trans)
}

func Translate(errs validator.ValidationErrors) error {
	var errList []string
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	errStr := strings.Join(errList, "|")
	return errors.New(errStr)
}

func Bind(ctx iris.Context, params interface{}, customize ...func(fl validator.FieldLevel) bool) (err error) {
	if ctx.Method() == "GET" {
		err = ctx.ReadQuery(params)
	} else {
		contentType := ctx.GetContentTypeRequested()
		if strings.Contains(contentType, "application/json") {
			err = ctx.ReadJSON(params)
		} else if contentType == "application/binary" {
			body, _ := ctx.GetBody()
			err = proto.Unmarshal(body, params.(proto.Message))
		} else {
			err = ctx.ReadForm(params)
		}
	}
	if err != nil {
		err = myError.Warp(err, "绑定数据出错")
		return
	}
	paramsStr, _ := json.Marshal(params)
	_, _ = ctx.Values().Set("params", string(paramsStr))
	// 注册自定义验证
	for _, validateFunc := range customize {
		funcName := GetFunctionName(validateFunc, '/', '.')
		_ = Validate.RegisterValidation(funcName, validateFunc)
	}
	// 参数验证
	err = Validate.Struct(params)
	if err != nil {
		err = Translate(err.(validator.ValidationErrors))
		return
	}
	return
}

func GetFunctionName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
