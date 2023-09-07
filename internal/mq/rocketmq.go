package mq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type Interface interface {
	Consumer() rocketmq.PushConsumer
	Producer() rocketmq.Producer
	Shutdown() (err error)
	Subscribe() error
}

type ConsumerExecFunc func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

type ConsumerInterface interface {
	GetTopic() string
	GetSelector() consumer.MessageSelector
	Exec() ConsumerExecFunc
}

type RocketMQ struct {
	producerProvider rocketmq.Producer
	consumerProvider rocketmq.PushConsumer
	conf             Config
}

type Config struct {
	NameServers primitive.NamesrvAddr
	GroupId     string
	Retries     int
}

func New(conf Config) (*RocketMQ, error) {
	mqClass := &RocketMQ{
		conf: conf,
	}
	if err := mqClass.newProducer(); err != nil {
		return nil, err
	}
	if err := mqClass.newConsumer(); err != nil {
		return nil, err
	}
	return mqClass, nil
}

func (r *RocketMQ) newProducer() error {
	var err error
	r.producerProvider, err = rocketmq.NewProducer(
		producer.WithNameServer(r.conf.NameServers),
		producer.WithRetry(r.conf.Retries),
		producer.WithGroupName(r.conf.GroupId),
	)
	if err != nil {
		return err
	}
	if err = r.producerProvider.Start(); err != nil {
		return err
	}
	return nil
}

func (r *RocketMQ) newConsumer() error {
	var err error
	r.consumerProvider, err = rocketmq.NewPushConsumer(
		consumer.WithGroupName(r.conf.GroupId),
		consumer.WithNameServer(r.conf.NameServers),
	)
	if err != nil {
		return err
	}
	if err = r.consumerProvider.Start(); err != nil {
		return err
	}
	return nil
}

func (r *RocketMQ) Subscribe(consumers ...ConsumerInterface) error {
	for _, cons := range consumers {
		if err := r.consumerProvider.Subscribe(cons.GetTopic(), cons.GetSelector(), cons.Exec()); err != nil {
			return err
		}
	}
	return r.startConsumer()
}

func (r *RocketMQ) startConsumer() error {
	return r.consumerProvider.Start()
}

func (r *RocketMQ) Consumer() rocketmq.PushConsumer {
	return r.consumerProvider
}

func (r *RocketMQ) Producer() rocketmq.Producer {
	return r.producerProvider
}

func (r *RocketMQ) Shutdown() (err error) {
	if err = r.producerProvider.Shutdown(); err != nil {
		return err
	}
	return r.consumerProvider.Shutdown()
}
