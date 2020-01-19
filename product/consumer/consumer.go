package main

import (
	"fmt"
	"product/common"
	"product/rabbitmq"
	"product/repositories"
	"product/services"
)

func main(){
	db ,err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	product := repositories.NewProductManager("product",db)
	productService := services.NewProductService(product)

	order := repositories.NewOrderManagerRepository("order",db)
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("product")
	rabbitmqConsumeSimple.ConsumeSimple(orderService,productService)
}
