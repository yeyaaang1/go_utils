package aliSLS

import (
	"fmt"
	"strconv"
)

type gormSLSLogger struct {
	agent SLSAgent
}

func (sls *gormSLSLogger) GormSuccessCallback(sql string, source string, rows int, duration float64) {
	_ = sls.agent.SendCustomizeMap(map[string]string{
		"sql":      sql,
		"file":     source,
		"rows":     strconv.Itoa(rows),
		"duration": fmt.Sprintf("%.2f", duration),
		"type":     "sql",
	})
}

func (sls *gormSLSLogger) GormFailCallback(source string, msg ...interface{}) {
	_ = sls.agent.SendCustomizeMap(map[string]string{
		"file": source,
		"msg":  fmt.Sprintf("%v", msg),
		"type": "sql",
	})
}

func NewGormSLSLogger(agent SLSAgent) GormSLSLogger {
	return &gormSLSLogger{agent}
}

type GormSLSLogger interface {
	GormSuccessCallback(sql string, source string, rows int, duration float64)
	GormFailCallback(source string, msg ...interface{})
}
