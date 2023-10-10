package amqp

import (
	"skeleton/app/amqp/consumer"
	"skeleton/internal/mq"
)

type Amqp struct{}

func (*Amqp) InitConsumers() []mq.ConsumerHandler {
	return []mq.ConsumerHandler{
		&consumer.FooConsumer{},
	}
}
