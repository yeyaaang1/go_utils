package aliSLS

import (
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/golog"
	"go_utils/myError"
	"runtime"
	"time"
)

type SLSAgent interface {
	// 设置程序标志
	// SetProject(project string)
	// 设置容器标志
	// SetLogStore(logStore string)
	// 发送日志
	Send(Contents []*sls.LogContent) error
	// 发送json
	// SendJson()
	// 发送map
	SendMap(mapContents map[string]string) (err error)
	// 带超时的关闭
	Close(timeoutMs int64) error
	// 安全关闭
	SafeClose()
	// logger的handler
	// SLSHandler(l *golog.Log) bool
}

var onlyProducer myProducer

type myProducer struct {
	producer *producer.Producer
	project  string
	logStore string
	topic    string

	level golog.Level
}

func (pro *myProducer) SetProject(project string) {
	pro.project = project
}

func (pro *myProducer) SetLogStore(logStore string) {
	pro.logStore = logStore
}

func (pro *myProducer) Send(contents []*sls.LogContent) (err error) {
	err = pro.send(pro.getSource(), contents)
	return
}

func (pro *myProducer) getSource(skip ...int) string {
	tmpSkip := 2
	if len(skip) > 0 {
		tmpSkip += skip[0]
	}
	pc, file, line, _ := runtime.Caller(tmpSkip)
	return fmt.Sprintf("%s(%d): %s", file, line, runtime.FuncForPC(pc).Name())
}

func (pro *myProducer) send(source string, contents []*sls.LogContent) (err error) {
	err = pro.producer.SendLog(pro.project, pro.logStore, pro.topic, source, &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	})
	return
}

func (pro *myProducer) SendMap(mapContents map[string]string) (err error) {
	var contents []*sls.LogContent
	for key, value := range mapContents {
		// 将value转换为string类型
		contents = append(contents, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(value),
		})
	}
	err = pro.send(pro.getSource(), contents)
	return

}

func (pro *myProducer) Close(timeoutMs int64) error {
	return pro.producer.Close(timeoutMs)
}

func (pro *myProducer) SafeClose() {
	pro.producer.SafeClose()
}

func (pro *myProducer) slsHandler(l *golog.Log) bool {
	if l.Level > pro.level || l.Level <= golog.DisableLevel {
		// 如果级别不符合预设要求, 则不进行上传
		return false
	}
	var loglevel string
	if level, ok := golog.Levels[l.Level]; ok {
		loglevel = level.Name
	}
	contents := []*sls.LogContent{
		{
			Key:   proto.String("level"),
			Value: proto.String(loglevel),
		}, {
			Key:   proto.String("msg"),
			Value: proto.String(l.Message),
		},
	}
	// fmt.Println(string(debug.Stack()))
	_ = pro.send(pro.getSource(4), contents)
	return false
}

func GetSLSHandler(level golog.Level) (handler func(l *golog.Log) bool, err error) {
	if onlyProducer.producer == nil {
		err = myError.New("SLS服务未初始化")
		return
	}
	onlyProducer.level = level
	handler = onlyProducer.slsHandler
	return
}

func GetSLS(endpoint, keyID, keySecret, project, logStore, topic string) SLSAgent {
	if onlyProducer.producer == nil {
		producerConfig := producer.GetDefaultProducerConfig()
		producerConfig.Endpoint = endpoint
		producerConfig.MaxIoWorkerCount = int64(runtime.NumCPU())
		producerConfig.MaxBlockSec = 0
		producerConfig.AccessKeyID = keyID
		producerConfig.AccessKeySecret = keySecret
		producerConfig.LogCompress = true
		producerInstance := producer.InitProducer(producerConfig)
		producerInstance.Start()
		onlyProducer = myProducer{
			producer: producerInstance,
			project:  project,
			logStore: logStore,
			topic:    topic,
			level:    golog.InfoLevel,
		}
	}
	return &onlyProducer
}
