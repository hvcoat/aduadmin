package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"hd/internal/biz"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type StudentService struct {
	stdUC *biz.StudentUseCase
}

func NewStudentService(uc *biz.StudentUseCase) *StudentService {
	return &StudentService{
		stdUC: uc,
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

		number := req.Form.Get("number")
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

		err = s.stdUC.SaveTask(req.Context(), number, taskID, req.RemoteAddr, &biz.Task{})
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
	number := req.Form.Get("number")
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
