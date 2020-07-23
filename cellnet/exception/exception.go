package exception

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"reflect"
	"runtime/debug"
)

func defaultCallback(raw interface{}, queue cellnet.EventQueue) {
	fmt.Println(fmt.Sprintf("%v\n%s", raw, string(debug.Stack())))
	queue.StopLoop()
}

func SetCapturePanic(queue cellnet.EventQueue, callback func(raw interface{}, queue cellnet.EventQueue)) {
	// 通过反射修改队列中异常捕获函数
	queueValue := reflect.ValueOf(queue)
	if callback == nil {
		callback = defaultCallback
	}
	queueValue.MethodByName("SetCapturePanicNotify").
		Call([]reflect.Value{reflect.ValueOf(callback)})
	// 开启异常捕获
	queue.EnableCapturePanic(true)
}
