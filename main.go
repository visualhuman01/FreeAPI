package main

import (
	"github.com/kataras/iris"
	)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views",".html"))
	app.Get("/", func(ctx iris.Context) {
		db := MysqlOperate{ ConnStr:"root:123456@tcp(139.219.13.4:3306)/eca?charset=utf8mb4"}
		rows := db.QueryData("select * from eca_course_schedules")
		for i,row := range rows{
			println("row:",i)
			for _,col := range row{
				if col != nil {
					print(string(col.([]byte)),",")
				}else{
					print("null,")
				}
			}
			println()
		}
		ctx.ViewData("message","Hello world!")
		ctx.View("hello.html")})
	app.Run(iris.Addr(":8888"))
}
