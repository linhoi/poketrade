package datamodels

type Order struct {
	ID          int64 `sql:"ID"`
	UserId      int64 `sql:"userId"`
	ProductId   int64 `sql:"productId"`
	OrderStatus int64 `sql:"orderStatus"`
}

const (
	OrderWait    = iota
	OrderSuccess  //1
	OrderFailed   //2
)

