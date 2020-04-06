package RMQ

import (
	"fmt"
	"sync"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

/*
Callback：RMQ接收消息的回调函数
MsgBody：消息体
messageID： 消息ID
ack：消息处理完后的回调
*/
type Callback func(MsgBody []byte, messageID string, ack func(messageID, messageMD5 string, err error) error) (err error)

type RMQOpt struct {
	mChanReadQName  string //读队列名
	mChanWriteQName string //写队列名
	mAmqpURI        string //amqr地址
	mExchangeName   string //
	mRoutingKey     string //路由
	mReadCallBack   Callback
	mReadConn       *amqp.Connection
	mNotify         chan *amqp.Error
	mChannel        *amqp.Channel

	mDeliveryLock sync.Mutex
	deliveryMap   map[string]*amqp.Delivery

	//消息ID处理次数，超过多少次失败后直接失败, msgID,count
	rmqMsgIDList map[string]int
}

//如果存在错误，则输出
func (opt *RMQOpt) failOnError(err error, msg string) error {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
	}
	return err
}

/*
函数说明：初始化
call: 业务层的回调函数
*/
func (opt *RMQOpt) InitMQTopic(AmqpURT, ExchnageName, ReadQName, WriteQName, RoutKey string, call Callback) (err error) {
	opt.mExchangeName = ExchnageName
	opt.mChanReadQName = ReadQName
	opt.mChanWriteQName = WriteQName
	opt.mAmqpURI = AmqpURT
	opt.mReadCallBack = call
	opt.mRoutingKey = RoutKey
	if opt.mReadCallBack != nil {
		opt.deliveryMap = make(map[string]*amqp.Delivery)
		go opt.consumer()
	}
	log.Debug("InitURL:", opt.mAmqpURI)
	return nil
}

/*
函数说明：发布消息到消息队列
MsgBody: 发布的消息
*/
func (opt *RMQOpt) Publish(MsgBody []byte) (err error) {
	connection, err := amqp.Dial(opt.mAmqpURI)
	if err != nil {
		log.Error("err:", err)
		return err
	}
	defer connection.Close()

	//创建一个Channel
	log.Info("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		log.Error("err:", err)
		return err
	}
	defer channel.Close()

	log.Printf("got queue, declaring %q", opt.mChanWriteQName)

	//创建一个queue
	q, err := channel.QueueDeclare(
		opt.mChanWriteQName, // name
		true,                // durable 持久化设置
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Error("err:", err)
		return err
	}

	if len(MsgBody) > 1024 {
		log.Printf("declared queue, publishing %dB", len(MsgBody))
	} else {
		log.Printf("declared queue, publishing %dB body (%q)", len(MsgBody), MsgBody)
	}

	// Producer只能发送到exchange，它是不能直接发送到queue的。
	// 现在我们使用默认的exchange（名字是空字符）。这个默认的exchange允许我们发送给指定的queue。
	// routing_key就是指定的queue名字。
	err = channel.Publish(
		opt.mExchangeName, // exchange
		q.Name,            // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			DeliveryMode:    amqp.Persistent, //持久化设置
			ContentEncoding: "",
			Body:            MsgBody,
		})
	if err != nil {
		log.Error("err:", err)
		return err
	}
	return nil
}

func (opt *RMQOpt) connectConsumer() error {
	if opt.mReadConn != nil {
		opt.mReadConn.Close()
	}
	var err error
	opt.mReadConn, err = amqp.Dial(opt.mAmqpURI)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ:", err)
		return err
	}
	opt.mNotify = opt.mReadConn.NotifyClose(make(chan *amqp.Error))
	log.Info("connect amqp consumer success")
	return nil
}

func (opt *RMQOpt) reconnectConsumer() error {
	for {
		if opt.connectConsumer() == nil {
			break
		}
		time.Sleep(5 * time.Second)
		log.Info("amqp consumer Reconnecting....")
	}
	return nil
}

//该方法用于与WechatServer通信，需要用到路由
func (opt *RMQOpt) getNormalRMQMsgChan() <-chan amqp.Delivery {
	//创建一个Channel
	var err error
	opt.mChannel, err = opt.mReadConn.Channel()
	opt.failOnError(err, "Failed to open a channel")

	log.Info("got queue, declaring ", opt.mChanReadQName)

	//创建一个exchange
	err = opt.mChannel.ExchangeDeclare(
		opt.mExchangeName, // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // noWait
		nil,               // arguments
	)
	opt.failOnError(err, "Failed to declare a queue")

	//创建一个queue
	q, err := opt.mChannel.QueueDeclare(
		opt.mChanReadQName, // name
		true,               // durable 持久化设置
		false,              // delete when unused
		false,              // exclusive 当Consumer关闭连接时，这个queue要被deleted
		false,              // no-wait
		nil,                // arguments
	)
	opt.failOnError(err, "Failed to declare a queue")

	//绑定到exchange
	err = opt.mChannel.QueueBind(
		opt.mChanReadQName, // name of the queue
		opt.mRoutingKey,    // bindingKey
		opt.mExchangeName,  // sourceExchange
		false,              // noWait
		nil,                // arguments
	)
	opt.failOnError(err, "Failed to declare a queue")

	//每次只取一条消息
	//这里为为了公平分发消息
	err = opt.mChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	opt.failOnError(err, "Failed to set QoS")

	log.Info("Queue bound to Exchange, starting Consume")
	//订阅消息
	msgs, err := opt.mChannel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil)
	opt.failOnError(err, "Failed to register a consumer")
	return msgs
}

/*
函数说明：消费者读取消息队列的消息，最终调用初始化函数注册的回调函数
*/
func (opt *RMQOpt) consumer() {
	//建立连接
	log.Info("dialing ", opt.mAmqpURI)
	if opt.connectConsumer() != nil {
		return
	}

	defer func() {
		log.Info("consumer exit")
		opt.mReadConn.Close()
	}()

	msgs := opt.getNormalRMQMsgChan()
	for { //receive loop
		select { //check connection
		case err := <-opt.mNotify:
			//work with error
			log.Error("work with error:", err)
			opt.mReadConn.Close()
			opt.mChannel.Close()
			if opt.reconnectConsumer() != nil {
				log.Error("Reconnect amqp consumer failed, exit")
				panic("RabbitMQ server is disconnect!!!!!")
				break
			}
			log.Info("Reconnect amqp consumer success")
			msgs = opt.getNormalRMQMsgChan()
		case d := <-msgs:
			opt.mDeliveryLock.Lock()
			//log.Debug("info:", d.DeliveryTag, ",len:", len(opt.deliveryMap))
			opt.deliveryMap[fmt.Sprint(d.DeliveryTag)] = &d
			opt.mDeliveryLock.Unlock()
			go opt.mReadCallBack(d.Body, fmt.Sprint(d.DeliveryTag), opt.ack)
		}
	}
}

func (opt *RMQOpt) ack(messageID, mesageMD5 string, err error) error {
	if err != nil {
		opt.rmqMsgIDList[string(mesageMD5)]++
		if opt.rmqMsgIDList[string(mesageMD5)] >= 10 {
			log.Debug("超过处理次数10，丢弃")
			delete(opt.rmqMsgIDList, mesageMD5)
			err = nil
		}
		time.Sleep(10 * time.Second)
	}
	opt.mDeliveryLock.Lock()
	defer opt.mDeliveryLock.Unlock()
	d, ok := opt.deliveryMap[messageID]
	if !ok {
		log.Error("err not find messageID:", messageID)
		return err
	}
	if err == nil {
		d.Ack(false)
	} else {
		d.Nack(false, true)
	}
	delete(opt.deliveryMap, messageID)
	return err
}

//PublishTopic topic模式的
func (opt *RMQOpt) PublishTopic(MsgBody []byte) (err error) {
	connection, err := amqp.Dial(opt.mAmqpURI)
	opt.failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	//创建一个Channel
	log.Info("got Connection, getting Channel")
	opt.mChannel, err = connection.Channel()
	opt.failOnError(err, "Failed to open a channel")
	defer opt.mChannel.Close()

	log.Info("got queue, declaring %q", opt.mChanWriteQName)
	//创建一个exchange
	err = opt.mChannel.ExchangeDeclare(
		opt.mExchangeName, // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // noWait
		nil,               // arguments
	)
	opt.failOnError(err, "Failed to declare a queue")

	if len(MsgBody) > 1024 {
		log.Printf("declared queue, publishing %dB", len(MsgBody))
	} else {
		log.Printf("declared queue, publishing %dB body (%q)", len(MsgBody), MsgBody)
	}

	// Producer只能发送到exchange，它是不能直接发送到queue的。
	// 现在我们使用默认的exchange（名字是空字符）。这个默认的exchange允许我们发送给指定的queue。
	// routing_key就是指定的路由名字。
	err = opt.mChannel.Publish(
		opt.mExchangeName, // exchange
		opt.mRoutingKey,   // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			DeliveryMode:    amqp.Persistent, //持久化设置
			ContentEncoding: "",
			Body:            MsgBody,
		})
	opt.failOnError(err, "Failed to publish a message")
	return nil
}
