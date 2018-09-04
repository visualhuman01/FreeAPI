package common

import "github.com/kataras/iris"

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

var Stu_Config = Config{}
var DBSource_Config map[int]DbConfig
var APP = iris.New()
func (p *DbConfig)GetDbConnStr() string {
	return p.Uid + ":" + p.Pwd + "@tcp(" + p.Ipaddr + ":" + p.Port + ")/" + p.Database + "?charset=utf8mb4"
}