package aliSLS

import (
	"fmt"
	"gitee.com/super_step/go_utils/myError"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/golog"
	"runtime"
	"time"
)

type SLSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Project         string `json:"project"`
	LogStore        string `json:"log_store"`
	Topic           string `json:"topic"`
}

type SLSAgent interface {
	// 发送自定义日志
	SendCustomize(contents []*sls.LogContent) error
	// 发送自定义map
	SendCustomizeMap(mapContents map[string]string) error
	// 发送日志
	Send(contents []*sls.LogContent, skip ...int) (err error)
	// 发送map
	SendMap(mapContents map[string]string) (err error)
	// 带超时的关闭
	Close(timeoutMs int64) error
	// 安全关闭
	SafeClose()
}

var onlyProducer myProducer

type myProducer struct {
	producer *producer.Producer
	project  string
	logStore string
	topic    string
	source   string

	level golog.Level
}

func (pro *myProducer) SetProject(project string) {
	pro.project = project
}

func (pro *myProducer) SetLogStore(logStore string) {
	pro.logStore = logStore
}

func (pro *myProducer) Send(contents []*sls.LogContent, skip ...int) (err error) {
	err = pro.send(pro.getSource(skip...), contents)
	return
}

func (pro *myProducer) getSource(skip ...int) (resSkip int) {
	resSkip = 2
	if len(skip) > 0 {
		resSkip += skip[0]
	}
	return resSkip
}

func (pro *myProducer) SendCustomize(contents []*sls.LogContent) error {
	return pro._send(contents)
}

func (pro *myProducer) SendCustomizeMap(mapContents map[string]string) error {
	var contents []*sls.LogContent
	for key, value := range mapContents {
		// 将value转换为string类型
		contents = append(contents, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(value),
		})
	}
	return pro._send(contents)
}

func (pro *myProducer) addDefaultType(contents []*sls.LogContent) {
	keyType := proto.String("type")
	for _, content := range contents {
		if content.Key == keyType {
			return
		}
	}
	contents = append(contents, &sls.LogContent{
		Key:   keyType,
		Value: proto.String("default"),
	})
}

func (pro *myProducer) _send(contents []*sls.LogContent) (err error) {
	pro.addDefaultType(contents)
	err = pro.producer.SendLog(pro.project, pro.logStore, pro.topic, pro.source, &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	})
	return
}

func (pro *myProducer) send(skip int, contents []*sls.LogContent) (err error) {
	pc, file, line, _ := runtime.Caller(skip)
	contents = append(contents, &sls.LogContent{
		Key:   proto.String("file"),
		Value: proto.String(fmt.Sprintf("%s(%d)", file, line)),
	}, &sls.LogContent{
		Key:   proto.String("func"),
		Value: proto.String(runtime.FuncForPC(pc).Name()),
	})
	return pro._send(contents)
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

func GetSLS(endpoint, keyID, keySecret, project, logStore, topic, source string) SLSAgent {
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
			source:   source,
		}
	}
	return &onlyProducer
}
