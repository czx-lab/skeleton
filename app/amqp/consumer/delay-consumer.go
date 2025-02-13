package consumer

import (
	"fmt"
	"skeleton/internal/mq"

	"github.com/streadway/amqp"
)

type DelayConsumer struct {
	queue    string
	consume  string
	exchange string
}

func NewDelayConsumer() *DelayConsumer {
	return &DelayConsumer{
		queue:    "delay-queue",
		consume:  "delay-consumer",
		exchange: "delay-delayed-exchange",
	}
}

func (d *DelayConsumer) Option() mq.ConsumerOption {
	return mq.ConsumerOption{
		CommonOption: mq.CommonOption{
			Mode:      mq.RoutingMode,
			ExChange:  d.exchange,
			QueueName: d.queue,
			Durable:   true,
			Kind:      "x-delayed-message",
			Args: amqp.Table{
				"x-delayed-type": "fanout",
			},
		},
		Consumer: d.consume,
	}
}

func (d *DelayConsumer) Exec(msg <-chan amqp.Delivery) {
	for v := range msg {
		go func(delivery amqp.Delivery) {
			if err := d.handler(delivery.Body); err != nil {
				delivery.Reject(true)
			} else {
				delivery.Ack(false)
			}
		}(v)
	}
}

func (*DelayConsumer) handler(message []byte) error {
	fmt.Println(message)
	return nil
}
