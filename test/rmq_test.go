package test

import (
	"fmt"
	"skeleton/internal/mq"
	"testing"

	"github.com/streadway/amqp"
)

var rmq = mq.NewRabbitMq("amqp://guest:guest@127.0.0.1:5672/")

type RmqFooConsumer struct{}

func (*RmqFooConsumer) Option() mq.ConsumerOption {
	return mq.ConsumerOption{
		CommonOption: mq.CommonOption{
			Mode:       mq.SimpleMode,
			QueueName:  "foo",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		AutoAck: true,
		NoLocal: false,
	}
}

func (*RmqFooConsumer) Exec(msg <-chan amqp.Delivery) {
	for v := range msg {
		fmt.Printf("consumer one message:%v\n", string(v.Body))
	}
}

type RmqFooTwoConsumer struct{}

func (*RmqFooTwoConsumer) Option() mq.ConsumerOption {
	return mq.ConsumerOption{
		CommonOption: mq.CommonOption{
			Mode:       mq.SimpleMode,
			QueueName:  "foo-two",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		AutoAck: true,
		NoLocal: false,
	}
}

func (*RmqFooTwoConsumer) Exec(msg <-chan amqp.Delivery) {
	for v := range msg {
		fmt.Printf("consumer two message:%v\n", string(v.Body))
	}
}

func TestRmqProducer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	if err := rmq.Publish(mq.ProducerOption{
		CommonOption: mq.CommonOption{
			Mode:       mq.SimpleMode,
			QueueName:  "foo",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		Message: amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World one!"),
		},
		Mandatory: false,
		Immediate: false,
	}); err != nil {
		t.Error(err)
	}

	if err := rmq.Publish(mq.ProducerOption{
		CommonOption: mq.CommonOption{
			Mode:       mq.SimpleMode,
			QueueName:  "foo-two",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		Message: amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World two!"),
		},
		Mandatory: false,
		Immediate: false,
	}); err != nil {
		t.Error(err)
	}
	t.Log("success")
}

func TestRmqConsumer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	rmq.Consumers(&RmqFooConsumer{}, &RmqFooTwoConsumer{})
	<-(chan any)(nil)
}
