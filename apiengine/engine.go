package apiengine

import (
	"strconv"
	"../common"
)

type Api_Interface struct {
	Method  string
	Input   []Api_Input
	Operate map[int]Api_Operate
	Output  Api_Output
}

type Api_Input struct {
	Name string
}

type Api_Output struct {
	Name      string
	Type      int //1:val,2:array_val,3:obj,4:array_obj
	OperateId int
	Fild      string
	Children  []*Api_Output
	Parent    *Api_Output
}

func (p *Api_Input) GetSymbol() string {
	return "@" + p.Name
}

type Api_Operate struct {
	DBSource_Id int
	SqlFormat   string
}
type ApiEngine struct {
}

var ApiInterface map[string]Api_Interface

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
	testapi := Api_Interface{Method: "Post"}
	testparam := Api_Input{Name: "id"}
	testapi.Input = make([]Api_Input, 1)
	testapi.Input[0] = testparam
	testapi.Operate = make(map[int]Api_Operate)
	testopt1 := Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test1 where f1 in (@id)"}
	testapi.Operate[1] = testopt1
	testopt2 := Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test2"}
	testapi.Operate[2] = testopt2
	testoutput := Api_Output{Name: "#", Type: 4, OperateId: 1, Fild: ""}
	testoutput1 := Api_Output{Name: "fid", Type: 1, OperateId: -1, Fild: "f1"}
	testoutput1.Parent = &testoutput
	testoutput2 := Api_Output{Name: "fname", Type: 1, OperateId: -1, Fild: "f2"}
	testoutput2.Parent = &testoutput
	testoutput3 := Api_Output{Name: "ftime", Type: 1, OperateId: -1, Fild: "f3"}
	testoutput3.Parent = &testoutput

	testoutput_sub := Api_Output{Name: "test2", Type: 4, OperateId: 2, Fild: ""}
	testoutput1_sub := Api_Output{Name: "fid", Type: 1, OperateId: -1, Fild: "ff1"}
	testoutput1_sub.Parent = &testoutput_sub
	testoutput2_sub := Api_Output{Name: "fname", Type: 1, OperateId: -1, Fild: "ff2"}
	testoutput2_sub.Parent = &testoutput_sub
	testoutput3_sub := Api_Output{Name: "ftime", Type: 1, OperateId: -1, Fild: "ff3"}
	testoutput3_sub.Parent = &testoutput_sub
	testoutput_sub.Children = make([]*Api_Output, 3)
	testoutput_sub.Children[0] = &testoutput1_sub
	testoutput_sub.Children[1] = &testoutput2_sub
	testoutput_sub.Children[2] = &testoutput3_sub

	testoutput.Children = make([]*Api_Output, 4)
	testoutput.Children[0] = &testoutput1
	testoutput.Children[1] = &testoutput2
	testoutput.Children[2] = &testoutput3
	testoutput.Children[3] = &testoutput_sub
	testapi.Output = testoutput
	ApiInterface = make(map[string]Api_Interface)
	ApiInterface["test"] = testapi
}
