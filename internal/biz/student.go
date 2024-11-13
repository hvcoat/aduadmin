package biz

import (
	"context"
	"fmt"
	"strings"
)

type Student struct{}

type Task struct {
	// 作业ID
	ID int
	// 作业内容
	Content string
}

type StudentRepo interface {
	// 获取签到页面
	Sign(ctx context.Context, name string, number, ip string) (string, error)
	// 签到
	Login(ctx context.Context, name string, number, ip string) error

	// 根据ID获取作业
	GetTask(ctx context.Context, id int) (string, error)
	// 提交作业
	SaveTask(ctx context.Context, name, number, ip string, task *Task) error
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
