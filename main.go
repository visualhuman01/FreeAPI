package main

import (
	"github.com/kataras/iris"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	DB   DB_Config
	Port string `json:"port"`
}

type DB_Config struct {
	Dbtype   string `json:"dbtype"`
	Ipaddr   string `json:"ipaddr"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Uid      string `json:"uid"`
	Pwd      string `json:"pwd"`
}

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	app.Get("/", func(ctx iris.Context) {
		connstr := Stu_Config.DB.Uid + ":" + Stu_Config.DB.Pwd + "@tcp(" + Stu_Config.DB.Ipaddr + ":" + Stu_Config.DB.Port + ")/" + Stu_Config.DB.Database + "?charset=utf8mb4"
		db := MysqlOperate{ConnStr: connstr}
		rows := db.QueryData("select * from eca_course_schedules")
		for i, row := range rows {
			println("row:", i)
			for _, col := range row {
				if col != nil {
					print(string(col.([]byte)), ",")
				} else {
					print("null,")
				}
			}
			println()
		}
		ctx.ViewData("message", "Hello world!")
		ctx.View("hello.html")
	})
	app.Run(iris.Addr(Stu_Config.Port))
}

var Stu_Config = Config{}

func init() {
	JsonParse := NewJsonStruct()
	JsonParse.Load("./config.json", &Stu_Config)
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
