package main

import "github.com/kataras/iris/v12"

func main(){
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.HandleDir("public","./frontend/web/public")
	app.HandleDir("html","./frontend/web/htmlProductShow")
	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		)
}
