package controllers

import (
	"github.com/kataras/iris"
	"../model"
	"../common"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"../apiengine"
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
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "insert into dbconfig_source set " +
		"source_name='" + m.Source_name + "'," +
		"source_type=" + m.Source_type + "," +
		"source_ipaddr='" + m.Source_ipaddr + "'," +
		"source_port=" + m.Source_port + "," +
		"source_database='" + m.Source_database + "'," +
		"source_uid='" + m.Source_uid + "'," +
		"source_pwd='" + m.Source_pwd + "'," +
		"source_des='" + m.Source_des + "'," +
		"source_status=1,source_createtime=now()"
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

func GetDBSourceListController(ctx iris.Context) {
	m := model.DBSource_Param{}
	ctx.ReadJSON(&m)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from dbconfig_source"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		println(err.Error())
	}
	ctx.JSON(data)
}
func BuildDBController(ctx iris.Context) {
	res := model.Result_Data{Code: 200, Msg: "ok"}
	m := model.BuildDB_Param{}
	ctx.ReadJSON(&m)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	connstr_tmp := common.DBSource_Config[m.Source_id]
	db_tmp := common.MysqlOperate{DBtype: connstr_tmp.Dbtype, ConnStr: connstr_tmp.GetDbConnStr()}
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from dbconfig_table where source_id = " + strconv.Itoa(m.Source_id)
	data1, err := db.QueryData(sqlstr)
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	for _, v := range data1 {
		tabid := v["table_id"].(int)
		tabname := v["table_name"].(string)
		sqlstr = "select a.*,b.datatype_name,b.datatype_is_fixed from dbconfig_field a inner join dbconfig_datatype b " +
			"on a.datatype_id = b.datatype_id where a.table_id = " + strconv.Itoa(tabid)
		data2, err := db.QueryData(sqlstr)
		if err != nil {
			res.Code = 500
			res.Msg = err.Error()
			ctx.JSON(res)
			return
		}
		mysql_sql := common.MysqlSQL{}
		sqlstr = mysql_sql.DropTable(tabname)
		db_tmp.Exec(sqlstr)
		sqlstr = mysql_sql.CreateTable(tabname, data2)
		db_tmp.Exec(sqlstr)
	}
	ctx.JSON(res)
}
func BuildTableController(ctx iris.Context) {
	res := model.Result_Data{Code: 200, Msg: "ok"}
	m := model.BuildTable_Param{}
	ctx.ReadJSON(&m)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	connstr_tmp := common.DBSource_Config[m.Source_id]
	db_tmp := common.MysqlOperate{DBtype: connstr_tmp.Dbtype, ConnStr: connstr_tmp.GetDbConnStr()}
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from dbconfig_table where table_id = " + strconv.Itoa(m.Table_id)
	data1, err := db.QueryRow(sqlstr)
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	tabname := data1["table_name"].(string)
	sqlstr = "select a.*,b.datatype_name,b.datatype_is_fixed from dbconfig_field a inner join dbconfig_datatype b " +
		"on a.datatype_id = b.datatype_id where a.table_id = " + strconv.Itoa(m.Table_id)
	data2, err := db.QueryData(sqlstr)
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	mysql_sql := common.MysqlSQL{}
	sqlstr = mysql_sql.DropTable(tabname)
	db_tmp.Exec(sqlstr)
	sqlstr = mysql_sql.CreateTable(tabname, data2)
	db_tmp.Exec(sqlstr)
	ctx.JSON(res)
}
func APIController(ctx iris.Context) {
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Origin", "*")
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "x-requested-with")
	ctx.ResponseWriter().Header().Add("content-type", "application/json")
	//ctx.ResponseWriter().Header().Add("Access-Control-Allow-Credentials","true")
	//ctx.ResponseWriter().Header().Add("Set-Cookie","company_auth=8A86YEGQQsS%252BJMXstQTIpqB78LBCLf9uaPb%252BV62h%2f6bHHMUdMx4Ge0EGl4XY2W4IZb9PCWzQon%2fxHZf8XZM8Y30BJHEB4%2ffAitm%2fwB6laKR3lhpminXvrcL8L9p2b4z%2fAaXKyzK3mnhEC%2fi6i7hc3aqTNVoVc5w9nT2PE%252BfpHcf8EFsKIXYHXg%3d%3d; expires=Sat, 08-Sep-2018 08:23:03 GMT; path=/")
	res := model.Result_Data{Code: 200, Msg: "ok"}
	aid := ctx.Params().Get("aid")
	api := apiengine.Apiengine.ApiInterface[aid]
	rawData, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(rawData, &dat); err == nil {
		output_tmp, err := apiengine.Apiengine.ApiOperate(api, dat)
		if err != nil {
			res.Code = 500
			res.Msg = err.Error()
			ctx.JSON(res)
			return
		}
		switch api.Output.Type {
		case 3:
			_, _, output_data, _ := apiengine.Apiengine.ApiOutput(api.Output, output_tmp, 0)
			ctx.JSON(output_data)
			return
		case 4:
			_, _, _, output_data := apiengine.Apiengine.ApiOutput(api.Output, output_tmp, 0)
			ctx.JSON(output_data)
			return
		}
	} else {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	ctx.JSON(res)
}
