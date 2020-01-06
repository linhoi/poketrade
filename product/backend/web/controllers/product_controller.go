package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/common"
	"product/datamodels"
	"product/services"
	"strconv"
)

type ProductController struct {
	Ctx 	iris.Context
	ProductService services.IProductService
}

func (p *ProductController) GetAll() mvc.View {
	productList, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name : "product/view.html",
		Data:iris.Map{
			"productList": productList,
		},
	}
}

func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOption(TagName:"form"))
	err := dec.Decode(p.Ctx.Request().Form, product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.ProductService.UpdateProduct(product)
	if err != nil {

		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name:"/product/add.html",
	}
}

func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOption(TagName:"form"))
	err := dec.Decode(p.Ctx.Request().Form, product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err = p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetManager() mvc.View {
	id ,err := strconv.ParseInt(p.Ctx.URLParam("id"),10,16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name:"product/manager.html",
		Data: iris.Map{
			"product":product,
		},
	}
}

func (p *ProductController) GetDelete() {
	idStr := p.Ctx.URLParam("id")
	id ,err := strconv.ParseInt(idStr,10,16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	isOk := p.ProductService.DeleteProductByID(id)
	if isOk {
		p.Ctx.Application().Logger().Debug("Delete the product success,ID is "+idStr )
	}else{
		p.Ctx.Application().Logger().Debug("Delete the product fail, ID is "+idStr )
	}

	p.Ctx.Redirect("/product/all")
}