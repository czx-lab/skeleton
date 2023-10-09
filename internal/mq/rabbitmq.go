package mq

import (
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	addr string
	pool sync.Pool
}

type RabbitMQConnect struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type RabbitMQInterface interface {
	Publish(opts ProducerOption) error
	Consume(opts ConsumerOption, handler ConsumerHandler) error
}

type ConsumerHandler interface {
	Exec(msg <-chan amqp.Delivery)
}

type Mode int

const (
	SimpleMode Mode = iota
	WorkMode
	SubscribeMode
	RoutingMode
	TopicMode
)

var modes = []string{"simple", "work", "subscribe", "routing", "topic"}

func (m Mode) String() string {
	return modes[m]
}

type CommonOption struct {
	Args amqp.Table

	// mq模式，默认simple
	Mode Mode

	// 持久化
	Durable bool

	// 队列中数据消费完成后是否自动删除队列
	AutoDelete bool

	// 是否独立， 为true时，适用于一个队列只能一个消费者
	Exclusive bool

	// 是否是内置的，如果设置为 true , 则表示内置交换器，客户端程序无法直接发送消息到这个交换器中
	Internal bool

	// 是否等待服务器返回
	NoWait bool

	// 队列名称（如果mode = routing || topic， queueName充当producer的key使用）
	QueueName string

	// 交换器
	ExChange string
	Kind     string
}

type ProducerOption struct {
	CommonOption
	Message amqp.Publishing

	// 是否强制
	Mandatory bool

	// 是否立即发送
	Immediate bool
}

type ConsumerOption struct {
	CommonOption
	Consumer string

	// 自动确认应答
	AutoAck bool

	// 非本地化，设置为true则表示不能将同一个Connection中生产者发送的消息传送给这个Connection中的消费者
	NoLocal bool
}

func NewRabbitMq(addr string) RabbitMQInterface {
	return &RabbitMQ{
		addr: addr,
		pool: sync.Pool{
			New: func() any {
				conn, err := connect(addr)
				if err != nil {
					panic(err)
				}
				return conn
			},
		},
	}
}

func (r *RabbitMQ) GetMQConn() *RabbitMQConnect {
	rmq := r.pool.Get()
	if rmq == nil {
		conn, err := connect(r.addr)
		if err != nil {
			panic(err)
		}
		return conn
	}
	return rmq.(*RabbitMQConnect)
}

func (r *RabbitMQ) Publish(opts ProducerOption) error {
	rmq := r.GetMQConn()
	defer r.pool.Put(rmq)
	name, err := r.producerQueueExchange(rmq, opts.CommonOption)
	if err != nil {
		return err
	}
	return rmq.channel.Publish(
		opts.ExChange,
		name,
		opts.Mandatory,
		opts.Immediate,
		opts.Message,
	)
}

func (r *RabbitMQ) Consume(opts ConsumerOption, handler ConsumerHandler) error {
	rmq := r.GetMQConn()
	queueName, err := r.consumerQueueExchange(rmq, opts.CommonOption)
	if err != nil {
		return err
	}
	msgs, err := rmq.channel.Consume(
		queueName,
		opts.Consumer,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}
	go func() {
		defer r.pool.Put(rmq)
		handler.Exec(msgs)
	}()
	return nil
}

func (r *RabbitMQ) queueDeclare(rmq *RabbitMQConnect, opts CommonOption) (amqp.Queue, error) {
	var err error
	queue, err := rmq.channel.QueueDeclare(
		opts.QueueName,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, nil
}

func (r *RabbitMQ) exchangeDeclare(rmq *RabbitMQConnect, opts CommonOption) error {
	return rmq.channel.ExchangeDeclare(
		opts.ExChange,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
}

func (r *RabbitMQ) producerQueueExchange(rmq *RabbitMQConnect, opts CommonOption) (string, error) {
	var err error
	var queue amqp.Queue
	if opts.Mode >= SubscribeMode {
		err = r.exchangeDeclare(rmq, opts)
	} else {
		queue, err = r.queueDeclare(rmq, opts)
	}
	if err != nil {
		return "", err
	}
	if opts.Mode >= RoutingMode {
		return opts.QueueName, nil
	}
	return queue.Name, nil
}

func (r *RabbitMQ) consumerQueueExchange(rmq *RabbitMQConnect, opts CommonOption) (string, error) {
	if opts.Mode >= SubscribeMode {
		if err := r.exchangeDeclare(rmq, opts); err != nil {
			return "", err
		}
	}
	queue, err := r.queueDeclare(rmq, opts)
	if err != nil {
		return "", err
	}
	if opts.Mode >= SubscribeMode {
		if err := rmq.channel.QueueBind(
			queue.Name,
			opts.QueueName,
			opts.ExChange,
			opts.NoWait,
			opts.Args,
		); err != nil {
			return "", err
		}
	}
	return queue.Name, nil
}

func connect(addr string) (*RabbitMQConnect, error) {
	var err error
	rmq := &RabbitMQConnect{}
	rmq.conn, err = amqp.Dial(addr)
	if err != nil {
		return nil, err
	}
	rmq.channel, err = rmq.conn.Channel()
	if err != nil {
		return nil, err
	}
	return rmq, nil
}
