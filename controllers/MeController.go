package controllers

import (
	"github.com/kataras/iris"
	"FreeAPI/common"
)

func HelloController(ctx iris.Context) {
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{ConnStr: connstr}
	rows := db.QueryData("select * from eca_course_schedules")
	for i, row := range rows {
		println("row:", i)
		for _, col := range row {
			if col != nil {
				print(string(col.([]byte)), ",")
			} else {
				print("null,")
			}
		}
		println()
	}
	ctx.ViewData("message", "Hello world!")
	ctx.View("hello.html")
}

func TestController(ctx iris.Context) {

	str := "test!!!"
	v := make(map[string]interface{})
	ctx.ReadJSON(&v)
	ctx.JSON(str)
}

func Test123Controller(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		return
	}
	ctx.WriteString("test:" + string(id))
}
