package middleware

import (
	"fmt"
	"gitee.com/super_step/go_utils/appConfig"
	"gitee.com/super_step/go_utils/iris/mail"
	"gitee.com/super_step/go_utils/timeTool"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"runtime/debug"
	"strings"
)

type errorLog struct {
	Method     string
	Host       string
	RequestURI string
	UserAgent  string
	ClientIP   string
	DebugStack string
	Error      interface{}
}

type exceptionMiddleware struct {
	mode       string
	appName    string
	logger     *golog.Logger
	mailConfig *appConfig.Mail
}

func (middleware *exceptionMiddleware) handler(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			go middleware.sendMail(errorLog{
				Error:      err,
				DebugStack: string(debug.Stack()),
				Method:     ctx.Method(),
				Host:       ctx.Host(),
				RequestURI: ctx.Request().URL.RequestURI(),
				UserAgent:  ctx.Request().UserAgent(),
				ClientIP:   ctx.Request().RemoteAddr + "(" + ctx.GetHeader("X-Real-IP") + ")",
			})
			ctx.StopExecution()
		}
	}()
	ctx.Next()
}

func (middleware *exceptionMiddleware) sendMail(log errorLog) {
	middleware.logger.Error("UserAgent: ", log.UserAgent)
	middleware.logger.Errorf("err: %s", log.Error)
	middleware.logger.Error(log.DebugStack)

	if middleware.mode != "debug" {
		subject := fmt.Sprintf("【重要错误】%s 项目出错了！", middleware.appName)

		debugStack := ""
		for _, v := range strings.Split(log.DebugStack, "\n") {
			debugStack += v + "<br>"
		}

		body := strings.ReplaceAll(mail.Template, "{ErrorMsg}", fmt.Sprintf("%s", log.Error))
		body = strings.ReplaceAll(body, "{RequestTime}", timeTool.GetCurrentDate())
		body = strings.ReplaceAll(body, "{RequestURL}", log.Method+"  "+log.Host+log.RequestURI)
		body = strings.ReplaceAll(body, "{RequestUA}", log.UserAgent)
		body = strings.ReplaceAll(body, "{RequestIP}", log.ClientIP)
		body = strings.ReplaceAll(body, "{DebugStack}", debugStack)

		options := &mail.Options{
			MailHost: middleware.mailConfig.Host,
			MailPort: middleware.mailConfig.Port,
			MailUser: middleware.mailConfig.User,
			MailPass: middleware.mailConfig.Pass,
			MailTo:   middleware.mailConfig.To,
			Subject:  subject,
			Body:     body,
		}

		_ = mail.Send(options)
	}
}

func GetExceptionMiddleware(mode, appName string, logger *golog.Logger, mailConfig *appConfig.Mail) func(ctx iris.Context) {
	middleware := exceptionMiddleware{
		mode:       mode,
		appName:    appName,
		logger:     logger,
		mailConfig: mailConfig,
	}
	return middleware.handler
}
