package timestamp

import (
	"time"
)

func GetTomorrow() int64 {
	timeStr := time.Now().Format("2006-01-02")
	// 使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	tomorrow, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return tomorrow.AddDate(0, 0, 1).Unix()
}
