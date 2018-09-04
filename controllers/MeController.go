package controllers

import (
	"github.com/kataras/iris"
	"../model"
	"../common"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"fmt"
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
	data := db.QueryData(sqlstr)
	ctx.JSON(data)
}
func BuildDBController(ctx iris.Context) {
	res := model.Result_Data{Code: 200, Msg: "ok"}
	m := model.BuildDB_Param{}
	ctx.ReadJSON(&m)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	//读取数据源配置
	sqlstr := "select * from dbconfig_source where source_id = " + strconv.Itoa(m.Source_id)
	data := db.QueryRow(sqlstr)
	dbconfig := common.DbConfig{}
	dbconfig.Dbtype = "mysql"
	dbconfig.Ipaddr = data["source_ipaddr"].(string)
	dbconfig.Port = strconv.Itoa(data["source_port"].(int))
	dbconfig.Database = data["source_database"].(string)
	dbconfig.Uid = data["source_uid"].(string)
	dbconfig.Pwd = data["source_pwd"].(string)
	//创建数据源实例
	connstr_tmp := dbconfig.GetDbConnStr()
	db_tmp := common.MysqlOperate{DBtype: dbconfig.Dbtype, ConnStr: connstr_tmp}
	sqlstr = "select * from dbconfig_table where source_id = " + strconv.Itoa(m.Source_id)
	data1 := db.QueryData(sqlstr)
	for _,v := range data1{
		tabid := v["table_id"].(int)
		tabname := v["table_name"].(string)
		sqlstr = "select a.*,b.datatype_name,b.datatype_is_fixed from dbconfig_field a inner join dbconfig_datatype b " +
			"on a.datatype_id = b.datatype_id where a.table_id = " + strconv.Itoa(tabid)
		data2 := db.QueryData(sqlstr)
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
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from dbconfig_source where source_id = " + strconv.Itoa(m.Source_id)
	data := db.QueryRow(sqlstr)
	dbconfig := common.DbConfig{}
	dbconfig.Dbtype = "mysql"
	dbconfig.Ipaddr = data["source_ipaddr"].(string)
	dbconfig.Port = strconv.Itoa(data["source_port"].(int))
	dbconfig.Database = data["source_database"].(string)
	dbconfig.Uid = data["source_uid"].(string)
	dbconfig.Pwd = data["source_pwd"].(string)
	//创建数据源实例
	connstr_tmp := dbconfig.GetDbConnStr()
	db_tmp := common.MysqlOperate{DBtype: dbconfig.Dbtype, ConnStr: connstr_tmp}
	sqlstr = "select * from dbconfig_table where table_id = " + strconv.Itoa(m.Table_id)
	data1 := db.QueryRow(sqlstr)
	tabname := data1["table_name"].(string)
	sqlstr = "select a.*,b.datatype_name,b.datatype_is_fixed from dbconfig_field a inner join dbconfig_datatype b " +
		"on a.datatype_id = b.datatype_id where a.table_id = " + strconv.Itoa(m.Table_id)
	data2 := db.QueryData(sqlstr)
	mysql_sql := common.MysqlSQL{}
	sqlstr = mysql_sql.DropTable(tabname)
	db_tmp.Exec(sqlstr)
	sqlstr = mysql_sql.CreateTable(tabname, data2)
	db_tmp.Exec(sqlstr)
	ctx.JSON(res)
}
func APIController(ctx iris.Context)  {
	res := model.Result_Data{Code: 200, Msg: "ok"}
	aid := ctx.Params().Get("aid")
	println(aid)
	rawData, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(rawData, &dat); err == nil {
		fmt.Println(dat)
	} else {
		fmt.Println(err)
	}
	println(string(rawData))
	ctx.JSON(res)
}