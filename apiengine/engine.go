package apiengine

import (
	"strconv"
	"../common"
	"strings"
)

var Apiengine = ApiEngine{}

type ApiEngine struct {
	ApiInterface map[string]common.Api_Interface
}

func (p *ApiEngine) Init() {
	//p.test()
	p.loadInterface()
	p.loadInput()
	p.loadOperate()
	p.loadOutput()
	p.loadBDSource()
}
func (p *ApiEngine) loadInterface() {
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	//读取数据源配置
	sqlstr := "select * from apiconfig_interface"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		panic(err)
	}
	p.ApiInterface = make(map[string]common.Api_Interface)
	for _, v := range data {
		interface_config := common.Api_Interface{Id: v["interface_id"].(int), Method: v["interface_method"].(int)}
		interface_config.Input = make([]common.Api_Input, 0)
		interface_config.Operate = make(map[int]common.Api_Operate)
		p.ApiInterface[v["interface_name"].(string)] = interface_config
	}
}
func (p *ApiEngine) loadInput() {
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	//读取数据源配置
	sqlstr := "select * from apiconfig_input"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		panic(err)
	}
	for k, v := range p.ApiInterface {
		input_tmp := make([]common.Api_Input,0)
		for _, vv := range data {
			if v.Id == vv["interface_id"].(int) {
				input_config := common.Api_Input{Name: vv["input_name"].(string)}
				input_tmp = append(input_tmp, input_config)
			}
		}
		tmp := p.ApiInterface[k]
		tmp.Input = input_tmp
		p.ApiInterface[k] = tmp
	}
}
func (p *ApiEngine) loadOperate() {
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	//读取数据源配置
	sqlstr := "select * from apiconfig_operate order by operate_step"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		panic(err)
	}
	for _, v := range data {
		for _, vv := range p.ApiInterface {
			if vv.Id == v["interface_id"].(int) {
				operate_config := common.Api_Operate{DBSource_Id: v["source_id"].(int), SqlFormat: v["operate_sql"].(string)}
				vv.Operate[v["operate_step"].(int)] = operate_config
			}
		}
	}
}
func (p *ApiEngine) loadOutput() {
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	//读取数据源配置
	sqlstr := "select a.*,b.field_name from apiconfig_output a left join dbconfig_field b on a.field_id = b.field_id"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		panic(err)
	}
	sqlstr = "SELECT a.*,b.field_name field_name,c.field_name val_field FROM apiconfig_condition a " +
		"left join dbconfig_field b on a.field_id = b.field_id " +
		"left join dbconfig_field c on a.condition_val_fieldid = c.field_id"
	data_c, err_c := db.QueryData(sqlstr)
	if err_c != nil {
		panic(err_c)
	}
	for k, v := range p.ApiInterface {
		tmp_data := make([]map[string]interface{}, 0)
		for _, vv := range data {
			if v.Id == vv["interface_id"] {
				tmp_data = append(tmp_data, vv)
			}
		}
		output_data := p.iterOutput(tmp_data, nil, data_c)
		if len(output_data) > 0 {
			tmp := p.ApiInterface[k]
			tmp.Output = output_data[0]
			p.ApiInterface[k] = tmp
		}
	}
}
func (p *ApiEngine) iterOutput(data []map[string]interface{}, parent interface{}, condition []map[string]interface{}) []common.Api_Output {
	res_output := make([]common.Api_Output, 0)
	for _, v := range data {
		if v["output_parent"] == parent {
			tmp_output := common.Api_Output{}
			tmp_output.Name = v["output_name"].(string)
			tmp_output.Type = v["output_type"].(int)
			tmp_output.OperateId = v["operate_id"].(int)
			if v["field_name"] != nil {
				tmp_output.Field = v["field_name"].(string)
			}
			tmp_condition := make([]map[string]interface{}, 0)
			for _, vv := range condition {
				if vv["ouput_id"] == v["output_id"] {
					tmp_condition = append(tmp_condition, vv)
				}
			}
			if len(tmp_condition) > 0 {
				condition_data := p.iterCondition(tmp_condition, nil)
				tmp_output.Condition = condition_data[0]
			}
			res_tmp := p.iterOutput(data, v["output_id"], condition)
			tmp_output.Children = make([]common.Api_Output, 0)
			for _, vv := range res_tmp {
				tmp_output.Children = append(tmp_output.Children, vv)
			}
			res_output = append(res_output, tmp_output)
		}
	}
	return res_output
}
func (p *ApiEngine) iterCondition(data []map[string]interface{}, parent interface{}) []common.Api_Condition {
	condition_data := make([]common.Api_Condition, 0)
	for _, v := range data {
		if v["condition_parent"] == parent {
			tmp_condition := common.Api_Condition{}
			if v["condition_type"] != nil {
				tmp_condition.Type = v["condition_type"].(int)
			}
			if v["condition_grouptype"] != nil {
				tmp_condition.GroupType = v["condition_grouptype"].(int)
			}
			if v["field_name"] != nil {
				tmp_condition.Field = v["field_name"].(string)
			}
			if v["condition_operator"] != nil {
				tmp_condition.Operator = v["condition_operator"].(int)
			}
			if v["condition_val_type"] != nil {
				tmp_condition.ValType = v["condition_val_type"].(int)
			}
			if v["condition_val_operateId"] != nil {
				tmp_condition.ValOperateId = v["condition_val_operateId"].(int)
			}
			if v["val_field"] != nil {
				tmp_condition.ValField = v["val_field"].(string)
			}
			if v["condition_val_datatype"] != nil {
				tmp_condition.ValDataType = v["condition_val_datatype"].(int)
			}
			tmp_condition.Val = v["condition_type"]
			tmp_res := p.iterCondition(data, v["condition_id"])
			tmp_condition.Children = make([]*common.Api_Condition, 0)
			for _, vv := range tmp_res {
				tmp_condition.Children = append(tmp_condition.Children, &vv)
			}
			condition_data = append(condition_data, tmp_condition)
		}
	}
	return condition_data
}
func (p *ApiEngine) loadBDSource() {
	connstr := common.Stu_Config.DB.GetDbConnStr()
	db := common.MysqlOperate{DBtype: common.Stu_Config.DB.Dbtype, ConnStr: connstr}
	//读取数据源配置
	sqlstr := "select * from dbconfig_source"
	data, err := db.QueryData(sqlstr)
	if err != nil {
		panic(err)
	}
	common.DBSource_Config = make(map[int]common.DbConfig)
	for _, v := range data {
		dbconfig := common.DbConfig{}
		dbconfig.Dbtype = "mysql"
		dbconfig.Ipaddr = v["source_ipaddr"].(string)
		dbconfig.Port = strconv.Itoa(v["source_port"].(int))
		dbconfig.Database = v["source_database"].(string)
		dbconfig.Uid = v["source_uid"].(string)
		dbconfig.Pwd = v["source_pwd"].(string)
		common.DBSource_Config[v["source_id"].(int)] = dbconfig
	}
}

//func (p *ApiEngine) test() {
//	testapi := common.Api_Interface{Method: 1}
//	testparam := common.Api_Input{Name: "id"}
//	testapi.Input = make([]common.Api_Input, 1)
//	testapi.Input[0] = testparam
//	testapi.Operate = make(map[int]common.Api_Operate)
//	testopt1 := common.Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test1 where f1 in (@id)"}
//	testapi.Operate[1] = testopt1
//	testopt2 := common.Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test2"}
//	testapi.Operate[2] = testopt2
//	testoutput := common.Api_Output{Name: "#", Type: 4, OperateId: 1, Field: ""}
//	testoutput1 := common.Api_Output{Name: "fid", Type: 1, OperateId: 1, Field: "f1"}
//	testoutput2 := common.Api_Output{Name: "fname", Type: 1, OperateId: 1, Field: "f2"}
//	testoutput3 := common.Api_Output{Name: "ftime", Type: 1, OperateId: 1, Field: "f3"}
//
//	testoutput_sub := common.Api_Output{Name: "test2", Type: 2, OperateId: 2, Field: "ff1"}
//	testcondition := common.Api_Condition{}
//	testcondition.Type = 2
//	testcondition.Field = "ff2"
//	testcondition.Operator = 1
//	testcondition.ValType = 2
//	testcondition.ValOperateId = 1
//	testcondition.ValField = "f1"
//	testcondition.ValDataType = 1
//	testoutput_sub.Condition = &testcondition
//	//testoutput1_sub := Api_Output{Name: "fid", Type: 1, OperateId: 2, Fild: "ff1"}
//	//testoutput2_sub := Api_Output{Name: "fname", Type: 1, OperateId: 2, Fild: "ff2"}
//	//testoutput3_sub := Api_Output{Name: "ftime", Type: 1, OperateId: 2, Fild: "ff3"}
//	//testoutput_sub.Children = make([]*Api_Output, 3)
//	//testoutput_sub.Children[0] = &testoutput1_sub
//	//testoutput_sub.Children[1] = &testoutput2_sub
//	//testoutput_sub.Children[2] = &testoutput3_sub
//
//	testoutput.Children = make([]*common.Api_Output, 4)
//	testoutput.Children[0] = &testoutput1
//	testoutput.Children[1] = &testoutput2
//	testoutput.Children[2] = &testoutput3
//	testoutput.Children[3] = &testoutput_sub
//	testapi.Output = testoutput
//	p.ApiInterface = make(map[string]common.Api_Interface)
//	p.ApiInterface["test"] = testapi
//}
func (p *ApiEngine) ApiCondition(condition common.Api_Condition, data map[int][]map[string]interface{}, OperateId int, index int) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if condition.Type == 1 {
		for _, v := range condition.Children {
			res_data := p.ApiCondition(*v, data, OperateId, index)
			if res_data != nil {
				for _, vv := range res_data {
					res = append(res, vv)
				}
			}
		}
	} else {
		var val interface{}
		if condition.ValType == 1 {
			val = condition.Val
		} else {
			rd := data[condition.ValOperateId]
			val = rd[index][condition.ValField]
		}
		ld := data[OperateId]
		for k, v := range ld {
			ldf := v[condition.Field]
			switch condition.Operator {
			case 1:
				if ldf == val {
					res = append(res, ld[k])
				}
				break
			case 2:
				break
			case 3:
				break
			case 4:
				break
			case 5:
				break
			case 6:
				break
			}
		}
	}
	return res
}
func (p *ApiEngine) ApiOutput(output common.Api_Output, data map[int][]map[string]interface{}, index int) (interface{}, []interface{}, map[string]interface{}, []map[string]interface{}) {
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
		var data_tmp []map[string]interface{}
		if output.Condition == nil {
			data_tmp = data[output.OperateId]
		} else {
			data_tmp = p.ApiCondition(output.Condition.(common.Api_Condition), data, output.OperateId, index)
		}
		if data_tmp != nil {
			res_val = data_tmp[index][output.Field]
		}
		break
	case 2:
		var data_tmp []map[string]interface{}
		if output.Condition == nil {
			data_tmp = data[output.OperateId]
		} else {
			data_tmp = p.ApiCondition(output.Condition.(common.Api_Condition), data, output.OperateId, index)
		}
		res_valarray = make([]interface{}, len(data_tmp))
		for k, v := range data_tmp {
			res_valarray[k] = v[output.Field]
		}
		break
	case 3:
		res_obj = make(map[string]interface{})
		for _, v := range output.Children {
			res1, res2, res3, res4 := p.ApiOutput(v, data, 0)
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
				res1, res2, res3, res4 := p.ApiOutput(v, data, k)
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
func (p *ApiEngine) ApiOperate(api common.Api_Interface, dat map[string]interface{}) (map[int][]map[string]interface{}, error) {
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
				tmpstr = strings.TrimRight(tmpstr, ",")
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
