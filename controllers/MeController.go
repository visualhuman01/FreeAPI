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
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Origin","*")
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers","x-requested-with")
	ctx.ResponseWriter().Header().Add("content-type","application/json")
	//ctx.ResponseWriter().Header().Add("Access-Control-Allow-Credentials","true")
	//ctx.ResponseWriter().Header().Add("Set-Cookie","company_auth=8A86YEGQQsS%252BJMXstQTIpqB78LBCLf9uaPb%252BV62h%2f6bHHMUdMx4Ge0EGl4XY2W4IZb9PCWzQon%2fxHZf8XZM8Y30BJHEB4%2ffAitm%2fwB6laKR3lhpminXvrcL8L9p2b4z%2fAaXKyzK3mnhEC%2fi6i7hc3aqTNVoVc5w9nT2PE%252BfpHcf8EFsKIXYHXg%3d%3d; expires=Sat, 08-Sep-2018 08:23:03 GMT; path=/")
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
		case 3:
			_, _, output_data, _ := apiOutput(api.Output, output_tmp, 0)
			ctx.JSON(output_data)
			return
		case 4:
			_, _, _, output_data := apiOutput(api.Output, output_tmp, 0)
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
func apiOutput(output apiengine.Api_Output, data map[int][]map[string]interface{}, index int) (interface{}, []interface{}, map[string]interface{}, []map[string]interface{}) {
	var res_val interface{}
	var res_valarray []interface{}
	var res_obj map[string]interface{}
	var res_objarray []map[string]interface{}
	res_val = nil
	res_valarray = nil
	res_obj = nil
	res_objarray = nil
	switch output.Type {
	case 1:
		data_tmp := data[output.OperateId]
		res_val = data_tmp[index][output.Fild]
		break
	case 2:
		data_tmp := data[output.OperateId]
		res_valarray = make([]interface{}, len(data_tmp))
		for k,v := range data_tmp{
			res_valarray[k] = v[output.Fild]
		}
		break
	case 3:
		res_obj = make(map[string]interface{})
		for _, v := range output.Children {
			res1, res2, res3, res4 := apiOutput(*v, data, 0)
			switch v.Type {
			case 1:
				res_obj[v.Name] = res1
				break
			case 2:
				res_obj[v.Name] = res2
				break
			case 3:
				res_obj[v.Name] = res3
				break
			case 4:
				res_obj[v.Name] = res4
				break
			}
		}
		break
	case 4:
		data_tmp := data[output.OperateId]
		res_objarray = make([]map[string]interface{}, len(data_tmp))
		for k, _ := range data_tmp {
			res_objarray[k] = make(map[string]interface{})
		}
		for _, v := range output.Children {
			for k, _ := range res_objarray {
				res1, res2, res3, res4 := apiOutput(*v, data, k)
				switch v.Type {
				case 1:
					res_objarray[k][v.Name] = res1
					break
				case 2:
					res_objarray[k][v.Name] = res2
					break
				case 3:
					res_objarray[k][v.Name] = res3
					break
				case 4:
					res_objarray[k][v.Name] = res4
					break
				}
			}
		}
		break
	}
	return res_val, res_valarray, res_obj, res_objarray
}
func apiOperate(api apiengine.Api_Interface, dat map[string]interface{}) (map[int][]map[string]interface{}, error) {
	operate_output := make(map[int][]map[string]interface{})
	for k, v := range api.Operate {
		sqlstr := ""
		for _, vv := range api.Input {
			tmpdat := dat[vv.Name]
			t, str, strarr := getJsonVal(tmpdat)
			if (t == 1) {
				sqlstr = strings.Replace(v.SqlFormat, vv.GetSymbol(), str, -1)
			} else {
				tmpstr := ""
				for _, v := range strarr {
					tmpstr += v + ","
				}
				tmpstr = strings.TrimRight(tmpstr,",")
				sqlstr = strings.Replace(v.SqlFormat, vv.GetSymbol(), tmpstr, -1)
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
func getJsonVal(jsondata interface{}) (int, string, []string) {
	var t int
	var res_str string
	var res_strarray []string
	switch jsondata.(type) {
	case string:
		t = 1
		res_str = jsondata.(string)
		break
	case float64:
		t = 1
		res_str = strconv.FormatFloat(jsondata.(float64), 'f', -1, 64)
		break
	case []interface{}:
		t = 2
		tmp := jsondata.([]interface{})
		res_strarray = make([]string, len(tmp))
		for k, v := range tmp {
			_, s, _ := getJsonVal(v)
			res_strarray[k] = s
		}
	}
	return t, res_str, res_strarray
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
