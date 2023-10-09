package mq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type Interface interface {
	Consumer() rocketmq.PushConsumer
	Producer() rocketmq.Producer
	TransProducer(listener primitive.TransactionListener) (rocketmq.TransactionProducer, error)
	Shutdown() (err error)
	Subscribe(consumers ...ConsumerInterface) error
	SendMessage(msg *primitive.Message) error
	SendTransactionMessage(rmq rocketmq.TransactionProducer, msg *primitive.Message) error
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
	conf             *Config
}

type Config struct {
	NameServers       primitive.NamesrvAddr
	ProducerGroupName string
	ConsumerGroupName string
	Retries           int
}

type Option interface {
	apply(*Config)
}

type OptionFunc func(conf *Config)

func (f OptionFunc) apply(conf *Config) {
	f(conf)
}

func New(opts ...Option) (Interface, error) {
	conf := &Config{}
	for _, opt := range opts {
		opt.apply(conf)
	}
	defaultConfig(conf)
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

func (r *RocketMQ) Subscribe(consumers ...ConsumerInterface) error {
	for _, cons := range consumers {
		if err := r.consumerProvider.Subscribe(cons.GetTopic(), cons.GetSelector(), cons.Exec()); err != nil {
			return err
		}
	}
	return r.consumerProvider.Start()
}

func (r *RocketMQ) SendMessage(msg *primitive.Message) error {
	_, err := r.producerProvider.SendSync(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (r *RocketMQ) SendTransactionMessage(rmq rocketmq.TransactionProducer, msg *primitive.Message) error {
	_, err := rmq.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (r *RocketMQ) Consumer() rocketmq.PushConsumer {
	return r.consumerProvider
}

func (r *RocketMQ) Producer() rocketmq.Producer {
	return r.producerProvider
}

func (r *RocketMQ) TransProducer(listener primitive.TransactionListener) (rocketmq.TransactionProducer, error) {
	var err error
	groupName := fmt.Sprintf("Trans%s", r.conf.ProducerGroupName)
	rmq, err := rocketmq.NewTransactionProducer(
		listener,
		producer.WithNameServer(r.conf.NameServers),
		producer.WithGroupName(groupName),
	)
	if err != nil {
		return nil, err
	}
	if err := rmq.Start(); err != nil {
		return nil, err
	}
	return rmq, nil
}

func (r *RocketMQ) Shutdown() (err error) {
	if err = r.producerProvider.Shutdown(); err != nil {
		return err
	}
	return r.consumerProvider.Shutdown()
}

func (r *RocketMQ) newProducer() error {
	var err error
	r.producerProvider, err = rocketmq.NewProducer(
		producer.WithNameServer(r.conf.NameServers),
		producer.WithRetry(r.conf.Retries),
		producer.WithGroupName(r.conf.ProducerGroupName),
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
		consumer.WithGroupName(r.conf.ConsumerGroupName),
		consumer.WithNameServer(r.conf.NameServers),
	)
	return err
}

func defaultConfig(conf *Config) {
	if conf.ConsumerGroupName == "" {
		conf.ConsumerGroupName = "defaultConsumerGroup"
	}
	if conf.ProducerGroupName == "" {
		conf.ProducerGroupName = "defaultProducerGroup"
	}
}

func WithNameServers(servers primitive.NamesrvAddr) Option {
	return OptionFunc(func(conf *Config) {
		conf.NameServers = servers
	})
}

func WithProducerGroupName(GId string) Option {
	return OptionFunc(func(conf *Config) {
		conf.ProducerGroupName = GId
	})
}

func WithConsumerGroupName(GId string) Option {
	return OptionFunc(func(conf *Config) {
		conf.ConsumerGroupName = GId
	})
}

func WithRetries(retries int) Option {
	return OptionFunc(func(conf *Config) {
		conf.Retries = retries
	})
}
