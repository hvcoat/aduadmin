package biz

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type Student struct{}

type Task struct {
	// 作业ID
	ID string
	// 作业内容
	Content string
}

type StudentRepo interface {
	// 获取签到页面
	Sign(ctx context.Context, name string, number, ip string) (string, error)
	// 签到
	Login(ctx context.Context, name string, number, ip string) error

	// 获取签到汇总信息
	ListSigns(ctx context.Context, date, gid, step string) ([]*StuSign, error)

	// 根据ID获取作业
	GetTask(ctx context.Context, id int) (string, error)
	// 提交作业
	SaveTask(ctx context.Context, name, number, ip string, task *Task) error
	// 获取签到汇总信息
	ListTaskSummary(ctx context.Context, date, gid, step string) ([]*StuTask, error)
}

type StudentUseCase struct {
	repo StudentRepo
}

func NewStudentUseCase(repo StudentRepo) *StudentUseCase {
	return &StudentUseCase{repo: repo}
}

func (s *StudentUseCase) Login(ctx context.Context, name string, number, ip string) error {
	return s.repo.Login(ctx, name, number, ip)
}

func (s *StudentUseCase) Sign(ctx context.Context) (string, error) {
	return sign, nil
}

func (s *StudentUseCase) GetTask(ctx context.Context, id int) (string, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		return "", err
	}

	// TODO: 使用html模板
	tmp := strings.ReplaceAll(taskContent, "{{.Task}}", fmt.Sprintf("%d: %s", id, task))
	tmp = strings.ReplaceAll(tmp, "{{.TaskID}}", fmt.Sprint(id))
	return tmp, nil
}

func (s *StudentUseCase) SaveTask(ctx context.Context, name, number, ip string, task *Task) error {
	return s.repo.SaveTask(ctx, name, number, ip, task)
}

type StuSign struct {
	Name   string
	Number string
	IP     string
	Date   string
}

type SignSummary struct {
	Summary []*StuSign
}

func (s *StudentUseCase) ListSigns(ctx context.Context, date, gid, step string) (*SignSummary, error) {
	stuSigns, err := s.repo.ListSigns(ctx, date, gid, step)
	if err != nil {
		log.Errorf("list signs error: %v, date: %s, gid: %s", err, date, gid)
		return nil, err
	}

	return &SignSummary{
		Summary: stuSigns,
	}, nil
}

type StuTask struct {
	Name   string
	Number string
	TaskID string
	IP     string
	Date   string
}

func (s *StudentUseCase) ListTaskSummary(ctx context.Context, date, gid, step string) ([]*StuTask, error) {
	stuTasks, err := s.repo.ListTaskSummary(ctx, date, gid, step)
	if err != nil {
		log.Errorf("list task summary error: %v, date: %s, gid: %s", err, date, gid)
		return nil, err
	}

	return stuTasks, nil
}
