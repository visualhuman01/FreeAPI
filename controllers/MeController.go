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
func GetTableListController(ctx iris.Context) {
	m := make(map[string]interface{})
	ctx.ReadJSON(&m)
	_, str, _ := common.GetJsonVal(m["sid"])
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from dbconfig_table where source_id = " + str
	data, err := db.QueryData(sqlstr)
	if err != nil {
		println(err.Error())
	}
	ctx.JSON(data)
}
func GetFieldListController(ctx iris.Context) {
	m := make(map[string]interface{})
	ctx.ReadJSON(&m)
	_, str, _ := common.GetJsonVal(m["tid"])
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select a.*,b.datatype_name from dbconfig_field a left join dbconfig_datatype b on a.datatype_id = b.datatype_id where a.table_id = " + str
	data, err := db.QueryData(sqlstr)
	if err != nil {
		println(err.Error())
	}
	ctx.JSON(data)
}
func GetInterfaceListController(ctx iris.Context) {
	m := model.DBSource_Param{}
	ctx.ReadJSON(&m)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from apiconfig_interface"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		println(err.Error())
	}
	ctx.JSON(data)
}
func AddTableViewController(ctx iris.Context) {
	sid := ctx.Params().Get("sid")
	ctx.ViewData("source_id", sid)
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr := "select * from dbconfig_datatype"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		println(err.Error())
	}
	ctx.ViewData("datatype", data)
	ctx.View("AddTable.html")
}
func AddTableController(ctx iris.Context) {
	res := model.Result_Data{Code: 200, Msg: "ok"}
	m := model.Table_Param{}
	ctx.ReadJSON(&m)
	//先建表
	sid, err := strconv.Atoi(m.Source_id)
	if err != nil {
		res.Code = 500
		res.Msg = "Source_id error"
		ctx.JSON(res)
		return
	}
	connstr_tmp := common.DBSource_Config[sid]
	db_tmp := common.MysqlOperate{DBtype: connstr_tmp.Dbtype, ConnStr: connstr_tmp.GetDbConnStr()}
	mysql_sql := common.MysqlSQL{}
	sqlstr := mysql_sql.DropTable(m.Table_name)
	db_tmp.Exec(sqlstr)
	field := make([]map[string]interface{}, 0)
	for _, v := range m.Field {
		r := common.Struct2Map(v,true)
		field = append(field, r)
	}
	sqlstr = mysql_sql.CreateTable(m.Table_name, field)
	db_tmp.Exec(sqlstr)
	//再插库
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	sqlstr = "insert into dbconfig_table(source_id,table_name,table_des,table_status,table_createtime,table_buildtime)" +
		" values(" + m.Source_id + ",'" + m.Table_name + "','" + m.Table_des + "',1,now(),now())"
	id := db.InsertData(sqlstr)
	for _, v := range m.Field {
		sqlstr = "insert into dbconfig_field(table_id,field_name,datatype_id,field_len,field_default,field_pk,field_null,field_auto,field_unsigned,field_zero,field_status,field_createtime,field_updatetime)" +
			"values(" + strconv.FormatInt(id, 10) + ",'" +
			"" + v.Field_name + "'," +
			"" + v.Datatype_id + "," +
			"" + v.Field_len + "," +
			"'" + v.Field_default + "'," +
			"" + strconv.Itoa(v.Field_pk) + "," +
			"" + strconv.Itoa(v.Field_null) + "," +
			"" + strconv.Itoa(v.Field_auto) + "," +
			"" + strconv.Itoa(v.Field_unsigned) + "," +
			"" + strconv.Itoa(v.Field_zero) + "," +
			"1,now(),now())"
		db.Exec(sqlstr)
	}
	ctx.JSON(res)
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
		sqlstr = "select a.*,b.datatype_name,b.datatype_is_fixed,b.datatype_is_quotation_mark from dbconfig_field a inner join dbconfig_datatype b " +
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
		sqlstr = "update dbconfig_table set table_buildtime = now() where table_id = " + strconv.Itoa(tabid)
		db.Exec(sqlstr)
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
	sqlstr = "select a.*,b.datatype_name,b.datatype_is_fixed,b.datatype_is_quotation_mark from dbconfig_field a inner join dbconfig_datatype b " +
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
	sqlstr = "update dbconfig_table set table_buildtime = now() where table_id = " + strconv.Itoa(m.Table_id)
	db.Exec(sqlstr)
	ctx.JSON(res)
}
func APIController(ctx iris.Context) {
	res := model.Result_Data{Code: 200, Msg: "ok"}
	aid := ctx.Params().Get("aid")
	api := apiengine.Apiengine.ApiInterface[aid]
	if api.IsCrossdomain {
		ctx.ResponseWriter().Header().Add("Access-Control-Allow-Origin", "*")
		ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "x-requested-with")
		ctx.ResponseWriter().Header().Add("content-type", "application/json")
	}
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
