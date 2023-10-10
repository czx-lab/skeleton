package consumer

import (
	"fmt"
	"skeleton/internal/mq"

	"github.com/streadway/amqp"
)

type FooConsumer struct{}

func (*FooConsumer) Option() mq.ConsumerOption {
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

func (*FooConsumer) Exec(msg <-chan amqp.Delivery) {
	for v := range msg {
		fmt.Printf("consumer message:%v\n", string(v.Body))
	}
}
