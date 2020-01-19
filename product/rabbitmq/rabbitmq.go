package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"product/datamodels"
	"product/services"
	"sync"
)
//url :   amqp://username:password@rabbitserver_address:server_port/virtual_host
const MQURL = "amqp://hikari:123456@127.0.0.1:5672/poketrade"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	QueueName string
	Exchange string
	Key string
	Mqurl string
	sync.Mutex
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ{
	return &RabbitMQ{
		QueueName:queueName,
		Exchange:exchange,
		Key:key,
		Mqurl:MQURL,
	}
}

func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName,"","")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err,"failed to open a channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(r.QueueName,false,false,false,false,nil)
	if err != nil {
		return err
	}
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err !=nil {
		return  err
	}
	return nil
}

func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService,productService services.IProductService ){
	//1. apply for a queue
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	_ = r.channel.Qos(
		1,
		0,
		false,
	)

	//2. receive message
	msgs, err := r.channel.Consume(
		q.Name,
		//saperate different consumer
		"",
		//ack to RabbitMQ when queue is consumed
		true,
		false,
		//if noLocal ==true, message can't be sent to this connection itself
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//3. consume queue
	forever := make (chan bool)
	go func(){
		for d := range msgs{
			log.Printf("Receive a message: %s",d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body),message)
			if err != nil {
				fmt.Println(err)
			}

			_, err = orderService.InserOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			err = productService.SubProductNum(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}

			d.Ack(false)

		}
	}()
	log.Printf("Waiting for messages, press CTRL + S to exit")
	<- forever
}
