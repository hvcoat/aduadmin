package data

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"hd/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewStudentRepo,
)

// Data .
type Data struct {
	// TODO wrapped database client

	// 先假设只使用只有两个文件
	login *os.File
	task  *os.File

	// 异步写
	content     chan string
	taskContent chan string
	saveDone    chan bool
	taskDone    chan bool
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	fileName := time.Now().Format("2006-01-02")
	fileName = fmt.Sprintf("%d-%s", c.Group.Gid, fileName)

	f, err := os.OpenFile(fileName+"-Login.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("打开文件失败: %v, fileName:%s\n", err, fileName)
		return nil, cleanup, err
	}

	ft, err := os.OpenFile(fileName+"-Task.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("打开文件失败: %v, fileName:%s\n", err, fileName)
		return nil, cleanup, err
	}

	d := &Data{
		login:       f,
		task:        ft,
		content:     make(chan string, 100), // 支持同时签到100个同学
		taskContent: make(chan string, 100), // 支持同时100个同学提交作业
		saveDone:    make(chan bool),        // 保存完成
		taskDone:    make(chan bool),        // 保存完成
	}

	go d.saveSign(context.TODO())
	go d.saveTask(context.TODO())

	// TODO:先这么写
	cleanup = func() {
		log.NewHelper(logger).Info("closing the data resources")
		close(d.content)
		<-d.saveDone
		<-d.taskDone
		f.Close()
		ft.Close()
	}

	return d, cleanup, nil
}

func (d *Data) SaveLogin(ctx context.Context, fileName string, content string) {
	d.content <- content
}

func (d *Data) SaveTask(ctx context.Context, fileName string, content string) {
	d.taskContent <- content
}

func (d *Data) saveSign(_ context.Context) {
	for c := range d.content {
		d.login.WriteString(string(c))
		d.login.WriteString("\n")
	}
	close(d.saveDone)
}

func (d *Data) saveTask(_ context.Context) {
	for c := range d.taskContent {
		d.task.WriteString(string(c))
		d.task.WriteString("\n")
	}
	close(d.taskDone)
}
