package main

import (
	"github.com/kataras/iris"
	"FreeAPI/common"
	"FreeAPI/controllers"
)

func main() {
	common.APP.RegisterView(iris.HTML("./views", ".html"))
	common.APP.StaticWeb("/", "./static") // serve our custom javascript code
	common.APP.Get("/",controllers.HelloController)
	common.APP.Any("/test",controllers.TestController)
	common.APP.Any("/test/{id}", controllers.Test123Controller)
	common.APP.Run(iris.Addr(common.Stu_Config.Port))
}

func init() {
	JsonParse := common.NewJsonStruct()
	JsonParse.Load("./config.json", &common.Stu_Config)
}


