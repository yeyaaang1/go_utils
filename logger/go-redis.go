package logger

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redisext"
	"github.com/kataras/golog"
	"runtime"
)

type GoRedisLogger struct {
	redisext.OpenTelemetryHook
	callbackLevel int
}

func NewGoRedisLogger(level int) redis.Hook {
	if level == 0 {
		level = 4
	}
	return &GoRedisLogger{callbackLevel: level}
}

// func (cl *GoRedisLogger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
// 	return cl.OpenTelemetryHook.BeforeProcess(ctx, cmd)
// }

func (cl *GoRedisLogger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	_, file, line, _ := runtime.Caller(cl.callbackLevel)
	golog.Default.Debugf("%s\n%s(%d)", cmd.String(), file, line)
	return cl.OpenTelemetryHook.AfterProcess(ctx, cmd)
}

// func (cl *GoRedisLogger) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
// 	return cl.OpenTelemetryHook.BeforeProcessPipeline(ctx, cmds)
// }

func (cl *GoRedisLogger) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	_, file, line, _ := runtime.Caller(cl.callbackLevel)
	golog.Default.Debugf("redis pipeline start: %s(%d)", file, line)
	for _, cmd := range cmds {
		golog.Default.Debug(cmd.String())
	}
	golog.Default.Debug("redis pipeline end")
	return cl.OpenTelemetryHook.AfterProcessPipeline(ctx, cmds)
}
