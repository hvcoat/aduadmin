package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	transporthttp "github.com/go-kratos/kratos/v2/transport/http"

	"hd/internal/biz"
	"hd/internal/conf"
)

type StudentService struct {
	stdUC      *biz.StudentUseCase
	taskNumber int32
}

func NewStudentService(d *conf.Data, uc *biz.StudentUseCase) *StudentService {
	return &StudentService{
		stdUC:      uc,
		taskNumber: d.TaskNumber,
	}
}

func (s *StudentService) Sign(w http.ResponseWriter, req *http.Request) {
	sign, err := s.stdUC.Sign(req.Context())
	if err != nil {
		fmt.Fprint(w, "获取签到页面失败")
		return
	}

	fmt.Fprint(w, sign)
}

func (s *StudentService) SignNew(ctx transporthttp.Context) error {
	sign, err := s.stdUC.Sign(ctx)
	if err != nil {
		return err
	}

	ctx.Response().Write([]byte(sign))
	return nil
}

func (s *StudentService) Pre(ctx transporthttp.Context) error {
	type ss struct {
		Name   string `json:"name"`
		TaskID string `json:"task-id"`
	}

	ctx.Request().Header.Set("Content-Type", "application/json")
	var sq ss
	err := ctx.BindVars(&sq)
	if err != nil {
		return err
	}
	name := sq.Name
	if name == "" {
		return fmt.Errorf("name is empty")
	}
	name = fmt.Sprintf("./%s/%s", sq.TaskID, name)

	for _, suffix := range []string{".png", ".jpg", ".jpeg"} {
		payload, err := loadPic(name + suffix)
		if err == nil {
			ctx.Response().Header().Set("Content-Type", "image/png")
			_, err = ctx.Response().Write(payload)
			return err
		}
	}

	return fmt.Errorf("not found")
}

func (s *StudentService) Index(ctx transporthttp.Context) error {

	// 先暂时只能看当天情况列表
	now := time.Now()
	step := "am"
	stepCH := "上午"
	beginFirstClass := "08:40:00"
	beginSencondClass := "10:30:00"
	if now.After(time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.Local)) {
		step = "pm"
		stepCH = "下午"
		beginFirstClass = "14:30:00"
		beginSencondClass = "16:30:00"
	}

	date := now.Format("2006-01-02")
	beginFirstClass = fmt.Sprintf(`%s-%s`, date, beginFirstClass)
	beginSencondClass = fmt.Sprintf(`%s-%s`, date, beginSencondClass)

	// 1组第一节课
	// 1组第二节课
	// 2组第一节课
	// 2组第二节课
	prefBlack := `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`
	classHref := []string{
		fmt.Sprintf(
			`%s<a href="/list-signs?date=%s&gid=1&step=%s&start-class-time=%s">今天1组%s第一节课签到汇总</a>`,
			prefBlack, date, step, beginFirstClass, stepCH),
		fmt.Sprintf(
			`%s<a href="/list-signs?date=%s&gid=1&step=%s&start-class-time=%s">今天1组%s第二节课签到汇总</a>`,
			prefBlack, date, step, beginSencondClass, stepCH),
		fmt.Sprintf(
			`%s<a href="/list-signs?date=%s&gid=2&step=%s&start-class-time=%s">今天2组%s第一节课签到汇总</a>`,
			prefBlack, date, step, beginFirstClass, stepCH),
		fmt.Sprintf(
			`%s<a href="/list-signs?date=%s&gid=2&step=%s&start-class-time=%s">今天2组%s第二节课签到汇总</a>`,
			prefBlack, date, step, beginSencondClass, stepCH),
	}

	var taskList []string
	for i := 1; i <= int(s.taskNumber); i++ {
		taskList = append(taskList,
			fmt.Sprintf(`%s<a href="/task?task-id=%d">作业%d</a>`, prefBlack, i, i),
		)
	}

	// http://127.0.0.1:8000/list-tasks?date=2024-11-18&gid=2&step=am
	taskSum := []string{
		fmt.Sprintf(`%s<a href="/list-tasks?date=%s&gid=1&step=%s">1组%s作业汇总</a>`, prefBlack, date, step, stepCH),
		fmt.Sprintf(`%s<a href="/list-tasks?date=%s&gid=2&step=%s">2组%s作业汇总</a>`, prefBlack, date, step, stepCH),
	}

	indexHtml := `<!DOCTYPE html>
<html lang="en">
<title>邯郸应用技术职业学院24级人工智能1班专用系统</title>
<font size="6" color="red" align="center">邯郸应用技术职业学院24级人工智能1班专用系统</font>
<h1>签到</h1>
%s<a href="/sign">签到</a>
<br/>
<br/>
%s

<br/>

<h1>作业</h1>
%s
<br/>
<br/>
%s
<br/>

</html>
`

	fh := fmt.Sprintf(indexHtml,
		prefBlack,
		strings.Join(classHref, "<br/>"),
		strings.Join(taskSum, "<br/>"),
		strings.Join(taskList, "<br/>"),
	)
	ctx.Response().Write([]byte(fh))
	return nil
}

func loadPic(name string) ([]byte, error) {
	image, err := os.Open(name)
	if err != nil {
		log.Errorf("open image error: %v", err)
		return nil, err
	}

	defer image.Close()
	payload, err := io.ReadAll(image)
	if err != nil {
		log.Errorf("read image error: %v", err)
		return nil, err
	}
	return payload, nil
}

func (s *StudentService) saveTask2Risk(taskID, number, typ, ext string, reader io.Reader) error {
	dir := fmt.Sprintf("./%s", taskID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
		if err != nil {
			fmt.Printf("创建文件夹失败,err=%v,dir=%s", err, dir)
			return fmt.Errorf("创建文件夹失败")
		}
	}

	name := filepath.Join(dir, fmt.Sprintf("%s-%s%s", number, typ, ext))
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("创建文件失败,err=%v,name=%s", err, name)
		return fmt.Errorf("创建文件失败,err=%v", err)
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("读取文件失败")
	}

	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("保存文件内容失败")
	}

	return nil
}

func (s *StudentService) getAndSaveTask(req *http.Request, taskID, number, key string) error {
	// 获取上传的文件
	file, handler, err := req.FormFile(key)
	if err != nil {
		return fmt.Errorf("获取文件失败")
	}
	defer file.Close()
	ext := filepath.Ext(handler.Filename)

	return s.saveTask2Risk(taskID, number, key, ext, file)
}

func (s *StudentService) SubmitTask(w http.ResponseWriter, req *http.Request) {
	// 检查请求方法是否为POST
	if req.Method == "POST" {
		// 解析multipart/form-data请求
		err := req.ParseMultipartForm(10 << 20) // 10MB内存限制，可调整
		if err != nil {
			fmt.Fprint(w, "无法解析表单数据")
			return
		}

		name := req.Form.Get("name")
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Fprint(w, "姓名不能为空")
			return
		}
		number := req.Form.Get("number")
		number = strings.TrimSpace(number)
		taskID := req.Form.Get("task-id")
		if number == "" {
			fmt.Fprint(w, "学号不能为空")
			return
		}

		for key := range req.MultipartForm.File {
			err = s.getAndSaveTask(req, taskID, number, key)
			if err != nil {
				fmt.Fprint(w, err.Error())
				return
			}
		}
		// 从RemoteAddr中提取IP部分
		index := strings.Index(req.RemoteAddr, ":")
		if index == -1 {
			return
		}

		ip := req.RemoteAddr[:index]

		err = s.stdUC.SaveTask(req.Context(), name, number, ip, &biz.Task{ID: taskID})
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Fprint(w, "作业提交成功")
	} else {
		fmt.Fprint(w, "只接受POST方法")
	}
}

func (s *StudentService) GetTask(w http.ResponseWriter, req *http.Request) {
	taskID := req.URL.Query().Get("task-id")
	id, _ := strconv.Atoi(taskID)
	task, err := s.stdUC.GetTask(req.Context(), id)
	if err != nil {
		fmt.Fprint(w, "获取作业失败,err=", err)
		return
	}

	_, _ = fmt.Fprintf(w, task)
}

func (s *StudentService) Login(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprint(w, "解析表单失败")
		return
	}

	name := req.Form.Get("name")
	name = strings.TrimSpace(name)
	number := req.Form.Get("number")
	number = strings.TrimSpace(number)
	if name == "" || number == "" {
		fmt.Fprint(w, "姓名或学号不能为空")
		return
	}

	// 从RemoteAddr中提取IP部分
	index := strings.Index(req.RemoteAddr, ":")
	if index == -1 {
		return
	}

	ip := req.RemoteAddr[:index]

	err = s.stdUC.Login(req.Context(), name, number, ip)
	if err != nil {
		fmt.Fprint(w, "签到失败")
		return
	}

	fmt.Fprint(w, "签到成功")
}
