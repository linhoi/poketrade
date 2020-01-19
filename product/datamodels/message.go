package datamodels

type Message struct {
	ProductID int64
	UserID    int64
}

func NewMessage(userId, productId int64) *Message {
	return &Message{userId,productId}
}