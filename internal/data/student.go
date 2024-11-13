package data

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"hd/internal/biz"
)

type studentRepo struct {
	data *Data
	log  *log.Helper

	// 模拟数据库
	task map[string]string
}

// GetTask implements biz.StudentRepo.
func (s *studentRepo) GetTask(ctx context.Context, id int) (string, error) {
	t, ok := s.task[fmt.Sprint(id)]
	if ok {
		return t, nil
	}
	fmt.Println("number=", len(s.task))
	for k, v := range s.task {
		fmt.Println(k, ":", v)
	}

	return "", fmt.Errorf("task not found: %d", id)
}

// Login implements biz.StudentRepo.
func (s *studentRepo) Login(ctx context.Context, name string, number string, ip string) error {
	t := time.Now().Format("2006-01-02-15:04:05")
	s.data.SaveLogin(ctx, "login", fmt.Sprintf("name:%s\tnumber:%s\ttime:%s\tip:%s", name, number, t, ip))
	return nil
}

// SaveTask implements biz.StudentRepo.
func (s *studentRepo) SaveTask(ctx context.Context, name string, number string, ip string, task *biz.Task) error {
	t := time.Now().Format("2006-01-02-15:04:05")
	s.data.SaveTask(ctx, "task", fmt.Sprintf("name:%s\tnumber:%s\ttime:%s\tip:%s\tcontent:%s", name, number, t, ip, task.Content))
	return nil
}

// Sign implements biz.StudentRepo.
func (s *studentRepo) Sign(ctx context.Context, name string, number string, ip string) (string, error) {
	panic("unimplemented")
}

// NewGreeterRepo .
func NewStudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	t, err := tasks()
	if err != nil {
		panic(err)
	}

	return &studentRepo{
		data: data,
		task: t,
		log:  log.NewHelper(logger),
	}
}

func tasks() (map[string]string, error) {
	file, err := os.Open("task.csv")
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return nil, err
	}
	defer file.Close()

	t := make(map[string]string)

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("读取文件出错: %v\n", err)
				return nil, err
			}

			if line == "" {
				break
			}
		}

		index := strings.Index(line, " ")
		if index == -1 {
			fmt.Printf("数据格式不正确,line=%s\n", line)
			break
		} else {
			t[line[:index]] = strings.TrimSpace(line[index+1:])
		}
	}

	return t, nil
}
