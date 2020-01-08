package services

import (
	"product/datamodels"
	"product/repositories"
)

type IOrderService interface {
	GetOrderByID(orderId int64) (order *datamodels.Order,err error)
	GetAllOrder()( []*datamodels.Order,  error)
	InsertOrder(order *datamodels.Order) (orderId int64, err error)
	GetAllOrderInfo() (map[int]map[string]string,  error)
	DeleteOrderByID(orderId int64) bool
	UpdateOrder(order *datamodels.Order) (err error)
}

type OrderService struct{
	OrderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderByID(orderId int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderId)
}

func (o *OrderService) GetAllOrder() ( []*datamodels.Order,error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (orderId int64, err error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrderInfo() ( map[int]map[string]string,  error) {
   return  o.OrderRepository.SelectAllWithInfo()
}

func (o *OrderService) DeleteOrderByID(orderId int64) bool {
	return o.OrderRepository.Delete(orderId)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	return o.OrderRepository.Update(order)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService{
	return &OrderService{repository}
}