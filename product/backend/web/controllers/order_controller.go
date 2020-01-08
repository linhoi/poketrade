package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/services"
)

type OrderController struct{
	Ctx iris.Context
	OrderService services.IOrderService  //a  bug service.OrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("Get Order Information Fail")
	}
	return mvc.View{
		Name:"order/view.html",
		Data: iris.Map{
			"order":orderArray,
		},
	}
}
