package common

import (
	"github.com/kataras/iris"
	"strconv"
	"reflect"
	"strings"
)

type Config struct {
	DB   DbConfig
	Port string `json:"port"`
}

type DbConfig struct {
	Dbtype   string `json:"dbtype"`
	Ipaddr   string `json:"ipaddr"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Uid      string `json:"uid"`
	Pwd      string `json:"pwd"`
}
type Api_Interface struct {
	Id            int
	Method        int //1:any,2:post,3:get
	IsCrossdomain bool
	Input         []Api_Input
	Operate       map[int]Api_Operate
	Output        Api_Output
}

type Api_Input struct {
	Name string
}

type Api_Output struct {
	Name      string
	Type      int //1:val,2:array_val,3:obj,4:array_obj
	OperateId int
	Field     string
	Children  []Api_Output
	Condition interface{}
}
type Api_Condition struct {
	Type         int //1:guroup,2:node
	GroupType    int //1:and,2:or
	Field        string
	Operator     int //1:=,2:>,3:<,4:>=,5:<=,6:!=
	ValType      int //1:val,2:OperateId
	ValOperateId int
	ValField     string
	ValDataType  int //1:number,2:string
	Val          interface{}
	Children     []*Api_Condition
}

func (p *Api_Input) GetSymbol() string {
	return "@" + p.Name
}

type Api_Operate struct {
	DBSource_Id int
	SqlFormat   string
}

var Stu_Config = Config{}
var DBSource_Config map[int]DbConfig
var APP = iris.New()

func (p *DbConfig) GetDbConnStr() string {
	return p.Uid + ":" + p.Pwd + "@tcp(" + p.Ipaddr + ":" + p.Port + ")/" + p.Database + "?charset=utf8mb4"
}

func GetJsonVal(jsondata interface{}) (int, string, []string) {
	var t int
	var res_str string
	var res_strarray []string
	switch jsondata.(type) {
	case string:
		t = 1
		res_str = jsondata.(string)
		break
	case int:
		t = 1
		res_str = strconv.Itoa(jsondata.(int))
		break
	case int64:
		t = 1
		res_str = strconv.FormatInt(jsondata.(int64), 10)
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
			_, s, _ := GetJsonVal(v)
			res_strarray[k] = s
		}
	default:
		println(reflect.TypeOf(jsondata))
		break
	}
	return t, res_str, res_strarray
}
func Struct2Map(obj interface{}, isLower bool) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		if isLower {
			key = strings.ToLower(key)
		}
		data[key] = v.Field(i).Interface()
	}
	return data
}
