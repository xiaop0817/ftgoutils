package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/xiaop0817/ftgoutils/c"
	"log"
	"time"
)

const (
	ConsumerTag = ""
	lc          = c.LightGreen
)

var prefix = c.C(fmt.Sprintf("[%-10s]", "RabbitMQ"), c.LightGreen)

var MQClient = &Client{
	connCreated: false,
	reConn:      false,
}

type HandlerFunc func(message amqp.Delivery) error

type Client struct {
	Connection    *amqp.Connection
	Channel       *amqp.Channel
	connCreated   bool
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	closeChan     chan *amqp.Error
	reConn        bool
}

func (client *Client) Start(userName string, pwd string, host string, port int) error {
	h := fmt.Sprintf("amqp://%s:%s@%s:%d/vas", userName, pwd, host, port)
	log.Printf("%s %s", prefix, c.C(h, lc))

	defer func() {
		if err := recover(); err != nil && !client.connCreated {
			fmt.Println(err)
			fmt.Printf("%s %s", prefix, c.C("RabbitMQ连接失败,再次连接..............", lc))
			time.Sleep(time.Second * 1)
			go client.Start(userName, pwd, host, port)
		}
	}()

	conn, err := amqp.Dial(h)
	client.Connection = conn
	if err != nil {
		failOnError(err, c.C("连接RabbitMQ失败", lc))
	}

	ch, err := conn.Channel()
	client.Channel = ch
	if err != nil {
		failOnError(err, c.C("打开通道失败", lc))
	}

	client.connNotify = client.Connection.NotifyClose(make(chan *amqp.Error))
	client.closeChan = make(chan *amqp.Error, 1)
	client.channelNotify = client.Channel.NotifyClose(client.closeChan)

	client.connCreated = true

	if !client.reConn {
		go client.ReConnect(userName, pwd, host, port)
	}
	log.Printf("%s %s", prefix, c.C("RabbitMQ连接完成", lc))
	return nil
}

func (client *Client) ReConnect(userName string, pwd string, host string, port int) {
	log.Printf("%s %s", prefix, c.C("RabbitMQ重连监听启动...............", lc))
	client.reConn = true
	for {
		select {
		case err := <-client.connNotify:
			if err != nil {
				log.Printf("%s %s 连接错误: %s %s", prefix, lc, err, c.End)
			}
		case err := <-client.channelNotify:
			if err != nil {
				log.Printf("%s %s 通道错误,准备重连: %s %s", prefix, lc, err, c.End)
			}
			client.Channel.Cancel(ConsumerTag, true)
			if client.Connection != nil && !client.Connection.IsClosed() {
				client.Connection.Close()
			}
			time.Sleep(2 * time.Second)
			client.Start(userName, pwd, host, port)
		}
	}
}

// AddListener 添加一个消费者
func (client *Client) AddListener(queueName string, handlerFunc HandlerFunc, routing string, exchange string, handlerName string) {
	var err error
	_, err = client.Channel.QueueDeclare(
		queueName,
		true,  //是否持久化
		false, //是否为自动删除
		false, //是否具有排他性
		false, //是否阻塞
		nil,   //额外属性
	)

	client.Channel.QueueBind(queueName, routing, exchange, false, nil)

	// chan amqp.Delivery 一个传递[amqp.Delivery]类型的channel,并将channel中数据写入到messages
	var messages <-chan amqp.Delivery
	messages, err = client.Channel.Consume(
		queueName,   // queue
		ConsumerTag, // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	go func(messages <-chan amqp.Delivery) {
		for message := range messages {
			if err := handlerFunc(message); err == nil {
				err := message.Ack(true)
				failOnError(err, "Failed to ACK!")
			} else {
				log.Printf("%s [%s] Failed to handle a message:[%s]", prefix, handlerName, c.C(err.Error(), c.Red))
				err := message.Reject(true)
				failOnError(err, "Failed to Reject!")
			}
		}
	}(messages)
	log.Printf("%s %s(%s)", prefix, c.C("RabbitMQ启动成功,正在监听消息", lc), c.C(handlerName, c.LightCyan))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s %s %s: %s %s", prefix, c.Red, msg, err, c.End)
	}
}
