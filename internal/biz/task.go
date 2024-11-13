package biz

var taskContent = `
<!DOCTYPE html>
<html>

<body>
	<font color="red" size="6">习题{{.Task}}</font>
	<hr/>
	<br/>
  <form action="submit-task" method="post" enctype="multipart/form-data">
  	姓 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  名：<input type="text" name="name"/>
	<br/>
  	学 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  号：<input type="text" name="number"/>
	<br/>
	<br/>
    代&nbsp;&nbsp;码&nbsp;&nbsp;文&nbsp;&nbsp;件:&nbsp;&nbsp;<input type="file" name="code" style="text-align: left;"/>
	<br/>
	<br/>
    运行结果截图:&nbsp;&nbsp;<input type="file" name="result" style="text-align: left;"/>
	<br/>
	<br/>
	<input type="hidden" name="task-id" value="{{.TaskID}}"/>
    <input type="submit" value="提交作业">
  </form>
</body>

</html>
`
