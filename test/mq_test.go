package test

import (
	"context"
	"fmt"
	_ "skeleton/internal/bootstrap"
	"skeleton/internal/mq"
	"skeleton/internal/variable"
	"testing"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func TestMqProducer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	msg := &primitive.Message{
		Topic: "foo-topic",
		Body:  []byte("test message!"),
	}
	if err := variable.MQ.SendMessage(msg); err != nil {
		t.Errorf("send fail: %s", err)
		return
	}
	t.Log("success")
}

type FooConsumer struct{}

func (*FooConsumer) GetTopic() string {
	return "foo-topic"
}

func (*FooConsumer) GetSelector() consumer.MessageSelector {
	return consumer.MessageSelector{}
}

func (*FooConsumer) Exec() mq.ConsumerExecFunc {
	return func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("message:%v", me)
		return consumer.ConsumeSuccess, nil
	}
}

func TestMqConsumer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	go func() {
		if err := variable.MQ.Subscribe(&FooConsumer{}); err != nil {
			t.Errorf("fail:%s", err)
		}
	}()
	select {}
}
