package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"os"
	"path/filepath"
	"product/datamodels"
	"product/services"
	"strconv"
	"text/template"
)

type ProductController struct {
	Ctx 	iris.Context
	ProductService 	services.IProductService
	OrderService	services.OrderService
	session			*sessions.Session
}


func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	return mvc.View{
		Layout:"shared/productLayout.html",
		Name: "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "抢购失败！"
	//判断商品数量是否满足需求
	if product.ProductNum > 0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		//创建订单
		userID, err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		order := &datamodels.Order{
			UserId:      int64(userID),
			ProductId:   int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		//新建订单
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功！"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}

}




//generate static html from template
var (
	//html store path
	htmlOutputPath = "./frontend/web/htmlProductShow/"
	//static template path
	templatePath = "./frontend/web/views/template/"
)
func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug()
	}

	contensTmp, err := template.ParseFiles(filepath.Join(templatePath,"product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	generateFileName := filepath.Join(htmlOutputPath,"htmlProduct.html")

	productData, err := p.ProductService.GetProductByID(int64(productID))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	generateStaticHtml(p.Ctx, contensTmp, generateFileName,productData)
}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	if exist(fileName){
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err !=nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file,product)
}
func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

