package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" form:"id"`
	ProductName  string `json:"productName" sql:"productName" form:"productName"`
	ProductNum   int64  `json:"productNum" sql:"productNum" form:"productNum"`
	ProductImage string `json:"productImage" sql:"productImage" form:"productImage"`
	ProductUrl   string `json:"productUrl" sql:"productUrl" form:"productUrl"`
}
