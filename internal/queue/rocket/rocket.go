package rocket

import (
	"errors"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var (
	NilProducerMsgErr = errors.New("producer message is empty")
)

type RocketConf struct {
	Servers   []string
	GroupName string
	Retry     int
}

type Rocket struct {
	conf RocketConf

	p rocketmq.Producer
	c rocketmq.PushConsumer
}

func New(conf RocketConf) *Rocket {
	return &Rocket{
		conf: conf,
		p:    producerInstance(conf),
		c:    consumerInstance(conf),
	}
}

func (r *Rocket) Producer() rocketmq.Producer {
	if r.p == nil {
		r.p = producerInstance(r.conf)
	}
	return r.p
}

func (r *Rocket) Consumer() rocketmq.PushConsumer {
	if r.c == nil {
		r.c = consumerInstance(r.conf)
	}
	return r.c
}

func consumerInstance(conf RocketConf) rocketmq.PushConsumer {
	options := []consumer.Option{
		consumer.WithNameServer(conf.Servers),
		consumer.WithGroupName(conf.GroupName),
		consumer.WithRetry(conf.Retry),
	}
	c, err := rocketmq.NewPushConsumer(options...)
	if err != nil {
		log.Fatalf("rocket consumer create error: %v", err.Error())
	}
	return c
}

func producerInstance(conf RocketConf) rocketmq.Producer {
	options := []producer.Option{
		producer.WithNameServer(conf.Servers),
		producer.WithGroupName(conf.GroupName),
		producer.WithRetry(conf.Retry),
	}
	p, err := rocketmq.NewProducer(options...)
	if err != nil {
		log.Fatalf("rocket producer create error: %v", err.Error())
	}
	if err := p.Start(); err != nil {
		log.Fatalf("rocket producer start error: %v", err.Error())
	}
	return p
}
