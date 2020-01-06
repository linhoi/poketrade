package main

import (
	"context"
	"database/sql"
	log "github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/backend/web/controllers"
	"product/common"
	"product/repositories"
	"product/services"
)

func main(){
	app := iris.New()
	app.Logger().SetLevel("gebug")
	template := iris.HTML("./backend/web/views",".hmtl").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	app.HandleDir("/assets","./backend/web/assets")
	app.OnAnyErrorCode(func(ctx iris.Context){
		ctx.ViewData("message","the page you are looking for is't exist")
		ctx.View("shared/error.html")
	})



	db ,err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	productRepository := repositories.NewProductManager("product", db)
	productService 	:= services.NewProductService(productRepository)
	productParty		:= app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))


	orderRepository := repositories.NewOrderManagerRepository("order",db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx,orderService)
	order.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(),
		iris.WithoutServerError(),
		)
}

func RegisterController(table string, db *sql.DB, party string, ctx context.Context){

}