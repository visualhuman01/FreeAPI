package main

import (
	"github.com/kataras/iris"
	"./common"
	"./controllers"
	"./apiengine"
)

func main() {
	common.APP.RegisterView(iris.HTML("./views", ".html"))
	common.APP.StaticWeb("/", "./static") // serve our custom javascript code
	RegisterRoute()
	common.APP.Run(iris.Addr(common.Stu_Config.Port))
}
func init() {
	JsonParse := common.NewJsonStruct()
	JsonParse.Load("./config.json", &common.Stu_Config)
	apiengine := apiengine.ApiEngine{}
	apiengine.Init()
}
func RegisterRoute()  {
	common.APP.Get("/",controllers.HelloController)
	common.APP.Get("/adddbsource", func(ctx iris.Context) {
		ctx.View("AddDBSource.html")
	})
	common.APP.Post("/adddbsource",controllers.AddDBSourceController)
	common.APP.Get("/dbsourceList", func(ctx iris.Context) {
		ctx.View("DBSourceList.html")
	})
	common.APP.Post("/dbsourceList",controllers.GetDBSourceListController)
	common.APP.Post("/builddb",controllers.BuildDBController)
	common.APP.Post("/buildtable",controllers.BuildTableController)
	common.APP.Post("/api/{aid}",controllers.APIController)
	//common.APP.Any("/test",controllers.TestController)
	//common.APP.Any("/test/{id}", controllers.Test123Controller)
}

