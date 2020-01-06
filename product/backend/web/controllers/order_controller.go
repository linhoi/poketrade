package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/services"
)

type OrderController struct{
	Ctx iris.Context
	OrderService services.OrderService
}

func (o *OrderController) GetAll() mvc.View {
	orderList, err := o.OrderService.GetAllOrder()
	if err != nil {
		o.Ctx.Application().Logger().Debug("Get Order Information Fail")
	}
	return mvc.View{
		Name:"order/view.html",
		Data: iris.Map{
			"orderList":orderList,
		},
	}
}
