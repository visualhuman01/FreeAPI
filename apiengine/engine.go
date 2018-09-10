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
	p.test()
	p.loadBDSource()
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
func (p *ApiEngine) test() {
	testapi := common.Api_Interface{Method: "Post"}
	testparam := common.Api_Input{Name: "id"}
	testapi.Input = make([]common.Api_Input, 1)
	testapi.Input[0] = testparam
	testapi.Operate = make(map[int]common.Api_Operate)
	testopt1 := common.Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test1 where f1 in (@id)"}
	testapi.Operate[1] = testopt1
	testopt2 := common.Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test2"}
	testapi.Operate[2] = testopt2
	testoutput := common.Api_Output{Name: "#", Type: 4, OperateId: 1, Fild: ""}
	testoutput1 := common.Api_Output{Name: "fid", Type: 1, OperateId: 1, Fild: "f1"}
	testoutput1.Parent = &testoutput
	testoutput2 := common.Api_Output{Name: "fname", Type: 1, OperateId: 1, Fild: "f2"}
	testoutput2.Parent = &testoutput
	testoutput3 := common.Api_Output{Name: "ftime", Type: 1, OperateId: 1, Fild: "f3"}
	testoutput3.Parent = &testoutput

	testoutput_sub := common.Api_Output{Name: "test2", Type: 2, OperateId: 2, Fild: "ff1"}
	testcondition := common.Api_Condition{}
	testcondition.Type = 2
	testcondition.Fild = "ff2"
	testcondition.Operator = 1
	testcondition.ValType = 2
	testcondition.ValOperateId = 1
	testcondition.ValFild = "f1"
	testcondition.ValDataType = 1
	testoutput_sub.Condition = &testcondition
	//testoutput1_sub := Api_Output{Name: "fid", Type: 1, OperateId: 2, Fild: "ff1"}
	//testoutput1_sub.Parent = &testoutput_sub
	//testoutput2_sub := Api_Output{Name: "fname", Type: 1, OperateId: 2, Fild: "ff2"}
	//testoutput2_sub.Parent = &testoutput_sub
	//testoutput3_sub := Api_Output{Name: "ftime", Type: 1, OperateId: 2, Fild: "ff3"}
	//testoutput3_sub.Parent = &testoutput_sub
	//testoutput_sub.Children = make([]*Api_Output, 3)
	//testoutput_sub.Children[0] = &testoutput1_sub
	//testoutput_sub.Children[1] = &testoutput2_sub
	//testoutput_sub.Children[2] = &testoutput3_sub

	testoutput.Children = make([]*common.Api_Output, 4)
	testoutput.Children[0] = &testoutput1
	testoutput.Children[1] = &testoutput2
	testoutput.Children[2] = &testoutput3
	testoutput.Children[3] = &testoutput_sub
	testapi.Output = testoutput
	p.ApiInterface = make(map[string]common.Api_Interface)
	p.ApiInterface["test"] = testapi
}
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
			val = rd[index][condition.ValFild]
		}
		ld := data[OperateId]
		for k, v := range ld {
			ldf := v[condition.Fild]
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
	return  res
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
		if output.Condition == nil{
			data_tmp = data[output.OperateId]
		}else{
			data_tmp = p.ApiCondition(*output.Condition,data,output.OperateId,index)
		}
		if data_tmp != nil{
			res_val = data_tmp[index][output.Fild]
		}
		break
	case 2:
		var data_tmp []map[string]interface{}
		if output.Condition == nil {
			data_tmp = data[output.OperateId]
		}else{
			data_tmp = p.ApiCondition(*output.Condition,data,output.OperateId,index)
		}
		res_valarray = make([]interface{}, len(data_tmp))
		for k, v := range data_tmp {
				res_valarray[k] = v[output.Fild]
		}
		break
	case 3:
		res_obj = make(map[string]interface{})
		for _, v := range output.Children {
			res1, res2, res3, res4 := p.ApiOutput(*v, data, 0)
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
				res1, res2, res3, res4 := p.ApiOutput(*v, data, k)
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
