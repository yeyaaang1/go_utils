package timeTool

import (
	"fmt"
	"time"
)

const (
	minute = 60
	hour   = 3600
	day    = 86400
)

// 获取当前的时间 - 字符串
func GetCurrentDate() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// 获取当前的时间 - Unix时间戳
func GetCurrentUnix() int64 {
	return time.Now().Unix()
}

// 获取当前的时间 - 毫秒级时间戳
func GetCurrentMilliUnix() int64 {
	return time.Now().UnixNano() / 1000000
}

// 获取当前的时间 - 纳秒级时间戳
func GetCurrentNanoUnix() int64 {
	return time.Now().UnixNano()
}

// 将时间转换为X秒, X分钟前, 这种形式
func TimeToDuration(inTime time.Time) string {
	timestamp := inTime.Unix()
	return TimeStampToDuration(timestamp)
}

// 将时间转换为X秒, X分钟前, 这种形式
func TimeStampToDuration(inTime int64) string {
	now := time.Now().Unix()
	seconds := now - inTime
	if seconds < 10 {
		return "刚刚"
	} else if seconds < minute {
		return fmt.Sprintf("%d秒前", seconds)
	} else if seconds < hour {
		return fmt.Sprintf("%d分钟前", int(seconds/minute))
	} else if seconds < day {
		return fmt.Sprintf("%d小时前", int(seconds/hour))
	} else {
		return fmt.Sprintf("%d天前", int(seconds/day))
	}
}
