package controllers

import (
	"github.com/kataras/iris"
	"../model"
	"../common"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"../apiengine"
	"strings"
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
	res := model.Result_Data{Code: 200, Msg: "ok"}
	aid := ctx.Params().Get("aid")
	api := apiengine.ApiInterface[aid]
	rawData, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(rawData, &dat); err == nil {
		output_tmp, err := apiOperate(api, dat)
		if err != nil {
			res.Code = 500
			res.Msg = err.Error()
			ctx.JSON(res)
			return
		}
		switch api.Output.Type {
		case 1:
			output_data := make(map[string]interface{})

			break
		case 2:
			break
		}
	} else {
		res.Code = 500
		res.Msg = err.Error()
		ctx.JSON(res)
		return
	}
	ctx.JSON(res)
}
func apiOutput(output apiengine.Api_Output,data map[int][]map[string]interface{})(int,map[string]interface{},
[]map[string]interface{},interface{})  {
	switch output.Type {
	case 1:
		output_data := make(map[string]interface{})
		for _,v := range output.Children{
			t,d1,d2,d3:=apiOutput(v,data)
			if t == 1{
				output[output.]
			}
		}
		break
	case 2:
		break
	case 3:
		dd := data[output.OperateId]
		val = dd
		return 3,nil,nil,
	}
}
func apiOperate(api apiengine.Api_Interface, dat map[string]interface{}) (map[int][]map[string]interface{}, error) {
	operate_output := make(map[int][]map[string]interface{})
	for k, v := range api.Operate {
		sqlstr := ""
		for _, vv := range api.Input {
			tmpdat := dat[vv.Name]
			switch tmpdat.(type) {
			case string:
				sqlstr = strings.Replace(v.SqlFormat, vv.GetSymbol(), tmpdat.(string), -1)
				break
			case float64:
				sqlstr = strings.Replace(v.SqlFormat, vv.GetSymbol(), strconv.FormatFloat(tmpdat.(float64), 'f', -1, 64), -1)
				break
			}
		}
		connstr_tmp := common.DBSource_Config[v.DBSource_Id]
		db_tmp := common.MysqlOperate{DBtype: connstr_tmp.Dbtype, ConnStr: connstr_tmp.GetDbConnStr()}
		data_tmp, err := db_tmp.QueryData(sqlstr)
		if err != nil {
			return nil, err
		}
		operate_output[k] = data_tmp
	}
	return operate_output, nil
}
func print_json(m interface{}) {
	switch vv := m.(type) {
	case string:
		fmt.Println(m, "is string", vv)
	case float64:
		fmt.Println(m, "is float", int64(vv))
	case int:
		fmt.Println(m, "is int", vv)
	case nil:
		fmt.Println(m, "is nil", "null")
	default:
		fmt.Println(m, "is of a type I don't know how to handle ", fmt.Sprintf("%T", m))
	}
}
