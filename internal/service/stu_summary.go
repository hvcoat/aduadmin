package service

import (
	"fmt"
	"hd/internal/biz"
	"hd/internal/conf"
	"html"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/samber/lo"
)

type StuSummaryService struct {
	stdUC *biz.StudentUseCase
	Studs []string
}

func NewStuSummaryService(d *conf.Data, uc *biz.StudentUseCase) *StuSummaryService {
	return &StuSummaryService{
		stdUC: uc,
		Studs: d.Group.Stus,
	}
}

func (stuS *StuSummaryService) ListSigns(w http.ResponseWriter, req *http.Request) {
	date := req.URL.Query().Get("date")
	if date == "" {
		fmt.Fprint(w, "日期不能为空")
		return
	}
	gid := req.URL.Query().Get("gid")
	if gid == "" {
		fmt.Fprint(w, "gid不能为空")
		return
	}

	step := req.URL.Query().Get("step")
	if step == "" {
		fmt.Fprint(w, "step不能为空")
		return
	}

	startClassTime := req.URL.Query().Get("start-class-time")
	if startClassTime == "" {
		fmt.Fprint(w, "start-class-time不能为空")
		return
	}

	signs, err := stuS.stdUC.ListSigns(req.Context(), date, gid, step)
	if err != nil {
		fmt.Fprintf(w, "list signs error: %v", err.Error())
		return
	}

	// table
	// | 姓名| 学号 | 签到时间 | 签到机器IP | 备注 |

	fmt.Fprint(w, "<!DOCTYPE html> <html><body>")
	fmt.Fprint(w, `<head>
    <style>
        table {
            border-collapse: collapse; /* 合并边框 */
        }

        table,
        th,
        td {
            border: 1px solid black; /* 给表格、表头、单元格都添加1px的黑色边框，这样就会显示出竖线 */
        }
    </style>
</head>`)
	showStep := "上午"
	if step == "pm" {
		showStep = "下午"
	}
	fmt.Fprintf(w, `<h1 align="center">%s %s 第%s组 签到情况</h1>`, date, showStep, gid)
	_, _ = fmt.Fprintf(w, `<table align="center"><tr><th>序号</th><th>姓名</th><th>学号</th><th>签到时间</th><th>签到机器IP</th><th>备注</th></tr>`)

	signStus := lo.Map(signs.Summary, func(item *biz.StuSign, index int) string {
		return item.Name
	})

	startTime, _ := time.Parse("2006-01-02-15:04:05", startClassTime)

	// left: unSign/moreSign
	left, right := lo.Difference(stuS.Studs, signStus)
	rightMap := lo.KeyBy(right, func(item string) string {
		return item
	})
	order := 1
	for _, s := range signs.Summary {
		comment := ""
		if _, ok := rightMap[s.Name]; ok {
			comment = "非本组学生"
		}
		d, _ := time.Parse("2006-01-02 15:04:05", s.Date)
		if d.After(startTime) {
			comment += `<font color="orange">迟到</font>`
		}

		_, _ = fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			order, s.Name, s.Number, s.Date, s.IP, comment)
		order++
	}

	for _, s := range left {
		comment := `<font color="red">旷课</font>`
		_, _ = fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			order, s, "", "", "", comment)
		order++
	}

	_, _ = fmt.Fprintf(w, "</table>")

	fmt.Fprint(w, "</body></html>")
}

func (stuS *StuSummaryService) ListTaskSummary(w http.ResponseWriter, req *http.Request) {
	date := req.URL.Query().Get("date")
	if date == "" {
		fmt.Fprint(w, "日期不能为空")
		return
	}
	gid := req.URL.Query().Get("gid")
	if gid == "" {
		fmt.Fprint(w, "gid不能为空")
		return
	}

	step := req.URL.Query().Get("step")
	if step == "" {
		fmt.Fprint(w, "step不能为空")
		return
	}

	tasks, err := stuS.stdUC.ListTaskSummary(req.Context(), date, gid, step)
	if err != nil {
		fmt.Fprintf(w, "list task summary error: %v", err.Error())
		return
	}

	// table
	// | 姓名| 学号 | 作业ID | 提交时间 | 签到机器IP | 作业内容 | 备注 |

	fmt.Fprint(w, "<!DOCTYPE html> <html><body>")
	fmt.Fprint(w, `<head>
    <style>
        table {
            border-collapse: collapse; /* 合并边框 */
        }

        table,
        th,
        td {
            border: 1px solid black; /* 给表格、表头、单元格都添加1px的黑色边框，这样就会显示出竖线 */
        }
		img {
  		    transition: all 0.5s ease;
    	}

	    img:hover {
 		    transform: scale(2.5);
    	}
    </style>
</head>`)
	showStep := "上午"
	if step == "pm" {
		showStep = "下午"
	}
	fmt.Fprintf(w, `<h1 align="center">%s %s 第%s组 作业提交情况</h1>`, date, showStep, gid)
	_, _ = fmt.Fprintf(w, `<table align="center">
	<tr>
		<th>序号</th>
		<th>姓名</th>
		<th>学号</th>
		<th>作业ID</th>
		<th>提交时间</th>
		<th>机器IP</th>
		<th>作业内容</th>
		<th>作业截图</th>
		<th>备注</th>
	</tr>`)

	taskStus := lo.Map(tasks, func(item *biz.StuTask, index int) string {
		return item.Name
	})

	// left: unSign/moreSign
	left, _ := lo.Difference(stuS.Studs, taskStus)
	// rightMap := lo.KeyBy(right, func(item string) string {
	// return item
	// })
	order := 1
	for _, s := range tasks {
		comment := ""
		// if _, ok := rightMap[s.Name]; ok {
		// 	comment = "非本组学生"
		// }

		_, _ = fmt.Fprintf(w, `<tr>
		<td>%d</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>
			<image src="%s"></image>
		</td>
		<td>%s</td>
		</tr>`,
			order,
			s.Name,
			s.Number,
			s.TaskID,
			s.Date,
			s.IP,
			fmt.Sprintf("<pre>%s</pre>",
				taskContent(fmt.Sprintf("./%s/%s-code", s.TaskID, s.Number)),
			), // 作业内容
			fmt.Sprintf("pre/%s/%s-result", s.TaskID, s.Number),
			comment,
		)
		order++
	}

	for _, s := range left {
		comment := `<font color="red">未提交作业</font>`
		_, _ = fmt.Fprintf(w, `<tr>
		<td>%d</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		</tr>`,
			order, s, "", "", "", "", "", "", comment)
		order++
	}

	_, _ = fmt.Fprintf(w, "</table>")

	fmt.Fprint(w, "</body></html>")
}

func taskContent(name string) string {
	var fileName string
	_, err := os.Stat(name + ".c")
	if err == nil {
		fileName = name + ".c"
	} else {
		_, err = os.Stat(name + ".cpp")
		if err != nil {
			return "未找到提交的作业文件"
		}
		fileName = name + ".cpp"
	}

	c, err := os.ReadFile(fileName)
	if err != nil {
		return "未找到提交的作业文件"
	}

	org := html.EscapeString(string(c))

	org = strings.ReplaceAll(org, "\n", "<br>")

	return org
}
