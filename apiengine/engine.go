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
	Name       string
	Type       int //1:obj,2:array,3:value
	OperateId  int
	ReturnName string
	Children   []Api_Output
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
	testopt := Api_Operate{DBSource_Id: 1, SqlFormat: "select * from test1 where f1=@id"}
	testapi.Operate = make(map[int]Api_Operate)
	testapi.Operate[1] = testopt
	testoutput := Api_Output{Name: "#", Type: 1, OperateId: 0, ReturnName: ""}
	testoutput1 := Api_Output{Name: "fid", Type: 3, OperateId: 1, ReturnName: "f1"}
	testoutput2 := Api_Output{Name: "fname", Type: 3, OperateId: 1, ReturnName: "f2"}
	testoutput3 := Api_Output{Name: "ftime", Type: 3, OperateId: 1, ReturnName: "f3"}
	testoutput.Children = make([]Api_Output, 3)
	testoutput.Children[0] = testoutput1
	testoutput.Children[1] = testoutput2
	testoutput.Children[2] = testoutput3
	testapi.Output = testoutput
	ApiInterface = make(map[string]Api_Interface)
	ApiInterface["test"] = testapi
}
