package repositories

import (
	"database/sql"
	"product/common"
	"product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() 	error
	Insert(order *datamodels.Order) (ID int64, err error)
	Delete(int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll () ([]*datamodels.Order, error)
	SelectAllWithInfo () (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql ,err := common.NewMysqlConn()
		if err != nil{
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == ""{
		o.table = "order"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productId int64, err error) {
	if err := o.Conn(); err != nil{
		return 0,err
	}
	sql := "insert into "+o.table+" userId=?,productId=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil{
		return 0,err
	}
	result, err := stmt.Exec(order.UserId,order.ProductId,order.OrderStatus)
	if err != nil{
		return 0,err
	}

	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderId int64) bool {
	if err := o.Conn(); err != nil{
		return false
	}
	sql := "delete from table "+o.table+" where ID = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil{
		return false
	}
	_, err = stmt.Exec(orderId)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil{
		return err
	}
	sql := "update order set userId=?,productId=?,orderStatus=? where ID = "+ strconv.FormatInt(order.ID, 10)
	stmt , err := o.mysqlConn.Prepare(sql)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(order.UserId,order.ProductId,order.OrderStatus)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderManagerRepository) SelectByKey(orderId int64) (orderResult *datamodels.Order, err error) {
	if err := o.Conn(); err != nil{
		return &datamodels.Order{}, err
	}
	sql := "select from table order where ID = " +strconv.FormatInt(orderId, 10)
	row , err := o.mysqlConn.Query(sql)
	if err != nil{
		return &datamodels.Order{},err
	}
	defer  row.Close()

	result := common.GetResultRow(row)
	if len(result) == 0 { return &datamodels.Order{}, err}
	common.DataToStructByTagSql(result, orderResult)
	return orderResult, nil

}

func (o *OrderManagerRepository) SelectAll() (orderList []*datamodels.Order,err error) {
	if err := o.Conn(); err != nil{
		return nil, err
	}
	sql := "select * from `order`"
	rows , err := o.mysqlConn.Query(sql)
	if err != nil{
		return nil,err
	}
	defer  rows.Close()


	result := common.GetResultRows(rows)
	if len(result) == 0 { return nil, err}
	for _, r := range result{
		oneOrder := &datamodels.Order{}
		common.DataToStructByTagSql(r, oneOrder)
		orderList = append(orderList,oneOrder)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	if err := o.Conn(); err != nil{
		return nil, err
	}
	sql := "select o.Id,p.productName,o.orderStatus from order as o left join product as p on o.productId=p.ID"
	rows , err := o.mysqlConn.Query(sql)
	if err != nil{
		return nil,err
	}
	defer  rows.Close()

	return  common.GetResultRows(rows), nil
}

func NewOrderManagerRepository(table string , mysqlConn *sql.DB) IOrderRepository{
	return &OrderManagerRepository{table,mysqlConn}
}