package middleware

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

type loggerMiddleware struct {
	logger *golog.Logger
}

func (middleware *loggerMiddleware) handler(ctx iris.Context) {
	ctx.Next()
	var logFunc func(format string, args ...interface{})
	record := ctx.Values().GetBoolDefault("record", false)
	statusCode := ctx.GetStatusCode()
	if statusCode >= 300 || record {
		if statusCode >= 300 {
			logFunc = middleware.logger.Warnf
		} else {
			logFunc = middleware.logger.Infof
		}
		logFunc("params: %s", ctx.Values().GetString("params"))
	}
	return
}

func GetLoggerMiddleware(logger *golog.Logger) func(ctx iris.Context) {
	middleware := loggerMiddleware{
		logger: logger,
	}
	return middleware.handler
}
