package controllers

import (
	"github.com/kataras/iris"
	"../model"
	"../common"
)

func HelloController(ctx iris.Context) {
	//connstr := common.Stu_Config.DB.GetDbConnStr()
	//db := common.MysqlOperate{ConnStr: connstr}
	//rows := db.QueryData("select * from eca_course_schedules")
	//for i, row := range rows {
	//	println("row:", i)
	//	for _, col := range row {
	//		if col != nil {
	//			print(string(col.([]byte)), ",")
	//		} else {
	//			print("null,")
	//		}
	//	}
	//	println()
	//}
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

func AddDBSourceController(ctx iris.Context) {
	m := model.DBSource_Param{}
	ctx.ReadJSON(&m)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{ConnStr: connstr}
	sqlstr := "insert into dbconfig_source set " +
		"dbconfig_source_name='" + m.Name + "'," +
		"dbconfig_source_ipaddr='" + m.Ipaddr + "'," +
		"dbconfig_source_port=" + m.Port + "," +
		"dbconfig_source_database='" + m.Database + "'," +
		"dbconfig_source_uid='" + m.Uid + "'," +
		"dbconfig_source_pwd='" + m.Pwd + "'"
	id := db.InsertData(sqlstr)
	res := model.Result_Data{}
	if id > 0 {
		res.Code = 200
		res.Msg = "ok"
	} else {
		res.Code = 500
		res.Msg = "error"
	}
	ctx.JSON(res)
}
