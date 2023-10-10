package producer

import (
	"skeleton/internal/mq"
	"skeleton/internal/variable"

	"github.com/streadway/amqp"
)

type FooProducer struct{}

func (*FooProducer) SendMessage(message []byte) error {
	opts := mq.ProducerOption{
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
			Body:        message,
		},
		Mandatory: false,
		Immediate: false,
	}
	if err := variable.Amqp.Publish(opts); err != nil {
		return err
	}
	return nil
}
