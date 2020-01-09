package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"product/common"
	"product/frontend/middleware"
	"product/frontend/web/controllers"
	"product/repositories"
	"product/services"
	"time"
)

func main(){
	app := iris.New()
	app.Logger().SetLevel("debug")
	template := iris.HTML("./frontend/web/views",".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	app.HandleDir("public","./frontend/web/public")

	app.HandleDir("html","./frontend/web/htmlProductShow")

	app.OnAnyErrorCode(func(ctx iris.Context){
		ctx.ViewData("message",ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {

	}
	sess := sessions.New(sessions.Config{
		Cookie:"AdminCookie",
		Expires:800*time.Minute,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user",db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx,sess.Start)
	userPro.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthBeforeLogin)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	_ = app.Run(
		iris.Addr("0.0.0.0:8081"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}