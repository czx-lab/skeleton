package producer

import (
	"encoding/json"
	"skeleton/internal/mq"
	"skeleton/internal/variable"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

// todo:: Requires rabbitmq delayed messaging plugin
// https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases
type DelayProducer struct {
	queue    string
	exchange string
}

var DelayedProducer = &DelayProducer{
	queue:    "delay-queue",
	exchange: "delay-delayed-exchange",
}

func (d *DelayProducer) Send(delay int, data any) error {
	message, err := json.Marshal(gin.H{"message": data})
	if err != nil {
		return err
	}
	_delay := time.Until(time.Unix(int64(delay), 0))
	opts := mq.ProducerOption{
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
		Message: amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
			Headers: amqp.Table{
				"x-delay": _delay.Milliseconds(),
			},
		},
	}
	if err := variable.Amqp.Publish(opts); err != nil {
		return err
	}
	return nil
}
