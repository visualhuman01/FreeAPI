package apiengine

import (
	"strconv"
	"../common"
)

type Api_Interface struct {
	Method  string
	Params  []Api_Params
	Operate []Api_Operate
}
type Api_Params struct {
	Name 		string
	SqlSymbol 	string
}
type Api_Operate struct {
	DBSource_Id	int
	SqlFormat 	string
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
	data := db.QueryData(sqlstr)
	common.DBSource_Config = make(map[int]common.DbConfig)
	for _,v := range data{
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
	testapi := Api_Interface{}
	testapi.Method = "Post"
	testparam := Api_Params{}
	testparam.Name = "id"
	testparam.SqlSymbol = "@id"
	testapi.Params = make([]Api_Params, 1)
	testapi.Params[0] = testparam
	testopt := Api_Operate{}
	testopt.DBSource_Id = 1
	testopt.SqlFormat = "select * from test1 f1=@id"
	testapi.Operate = make([]Api_Operate, 1)
	testapi.Operate[0] = testopt
	ApiInterface = make(map[string]Api_Interface)
	ApiInterface["test"] = testapi
}
