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
	s.data.SaveTask(ctx, "task", fmt.Sprintf("name:%s\tnumber:%s\ttask-id:%s\ttime:%s\tip:%s\tcontent:%s", name, number, task.ID, t, ip, task.Content))
	return nil
}

// Sign implements biz.StudentRepo.
func (s *studentRepo) Sign(ctx context.Context, name string, number string, ip string) (string, error) {
	panic("unimplemented")
}

func (s *studentRepo) ListSigns(ctx context.Context, date, gid, step string) ([]*biz.StuSign, error) {
	name := fmt.Sprintf("%s-%s-%s-Login.csv", gid, date, step)
	stuSigns, err := signs(name)
	if err != nil {
		return nil, err
	}

	return stuSigns, nil
}

func (s *studentRepo) ListTaskSummary(ctx context.Context, date, gid, step string) ([]*biz.StuTask, error) {
	name := fmt.Sprintf("%s-%s-%s-Task.csv", gid, date, step)
	stuTask, err := stuTasks(name)
	return stuTask, err
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

func stuTasks(name string) ([]*biz.StuTask, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("打开文件失败: %v, name:%s\n", name, err)
		return nil, err
	}
	defer file.Close()

	var ret []*biz.StuTask
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

		// 去除换行
		line = line[:len(line)-1]
		// ("name:%s\tnumber:%s\ttask-id:%d\ttime:%s\tip:%s\tcontent:%s", name, number, t, ip, task.Content))
		datas := strings.Split(line, "\t")
		if len(datas) != 6 {
			log.Errorf("数据格式不正确,line=%s", line)
			return nil, fmt.Errorf("数据格式不正确,line=%s", line)
		}

		date := segValue(datas[3])
		lIndex := strings.LastIndex(date, "-")
		date = date[:lIndex] + " " + date[lIndex+1:]

		sign := &biz.StuTask{
			Name:   segValue(datas[0]),
			Number: segValue(datas[1]),
			TaskID: segValue(datas[2]),
			Date:   date,
			IP:     segValue(datas[4]),
		}

		ret = append(ret, sign)
	}

	return ret, nil
}

func signs(name string) ([]*biz.StuSign, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("打开文件失败: %v, name:%s\n", name, err)
		return nil, err
	}
	defer file.Close()

	var ret []*biz.StuSign
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

		line = line[:len(line)-1]
		datas := strings.Split(line, "\t")
		if len(datas) != 4 {
			log.Errorf("数据格式不正确,line=%s", line)
			return nil, fmt.Errorf("数据格式不正确,line=%s", line)
		}

		date := segValue(datas[2])
		lIndex := strings.LastIndex(date, "-")
		date = date[:lIndex] + " " + date[lIndex+1:]

		sign := &biz.StuSign{
			Name:   segValue(datas[0]),
			Number: segValue(datas[1]),
			Date:   date,
			IP:     segValue(datas[3]),
		}

		ret = append(ret, sign)
	}

	return ret, nil
}

func segValue(value string) string {
	index := strings.Index(value, ":")
	if index == -1 {
		return value
	}

	return value[index+1:]
}
