package logger

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

type MyDBLogger struct {
	SuccessCallback func(sql string, source string, rows int, duration float64)
	FailCallback    func(source string, msg ...interface{})
}

func (logger *MyDBLogger) Print(values ...interface{}) {
	// fmt.Println(fmt.Sprintf("%+v", values))
	var (
		level       = values[0]
		currentTime = "\n\033[33m[" + time.Now().Format("2006-01-02 15:04:05") + "]\033[0m"
		messages    = []interface{}{fmt.Sprintf("\033[35m(%v)\033[0m", values[1]), currentTime}
	)
	if level == "sql" {
		// fmt.Println(sql, source, values[5], values[2])
		// 构建原始的控制台输出
		duration := float64(values[2].(time.Duration).Nanoseconds()/1e4) / 100.0
		messages = append(messages, fmt.Sprintf(" \033[36;1m[%.2fms]\033[0m ", duration))
		sql := sqrFormat(values[3].(string), values[4].([]interface{}))
		messages = append(messages, sql)
		messages = append(messages, fmt.Sprintf(" \n\033[36;31m[%v]\033[0m ", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned "))
		if logger.SuccessCallback != nil {
			rows := values[5].(int)
			source := fmt.Sprintf("%v", values[1])
			logger.SuccessCallback(sql, source, rows, duration)
		}
	} else {
		messages = append(messages, "\033[31;1m")
		messages = append(messages, values[2:]...)
		messages = append(messages, "\033[0m")
		if logger.FailCallback != nil {
			source := fmt.Sprintf("%v", values[1])
			logger.FailCallback(source, values[2:]...)
		}
	}
	fmt.Println(messages...)
}

func sqrFormat(inSql string, args []interface{}) (sql string) {
	var (
		formattedValues []string
	)
	for _, value := range args {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() {
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok {
				if t.IsZero() {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
				}
			} else if b, ok := value.([]byte); ok {
				if str := string(b); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			} else if r, ok := value.(driver.Valuer); ok {
				if value, err := r.Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			} else {
				switch value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
					formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
				default:
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
		} else {
			formattedValues = append(formattedValues, "NULL")
		}
	}

	// differentiate between $n placeholders or else treat like ?
	if numericPlaceHolderRegexp.MatchString(inSql) {
		sql = inSql
		for index, value := range formattedValues {
			placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
			sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
		}
	} else {
		formattedValuesLength := len(formattedValues)
		for index, value := range sqlRegexp.Split(inSql, -1) {
			sql += value
			if index < formattedValuesLength {
				sql += formattedValues[index]
			}
		}
	}
	return
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
