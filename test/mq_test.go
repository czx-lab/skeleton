package test

import (
	"context"
	"fmt"
	_ "skeleton/internal/bootstrap"
	"skeleton/internal/mq"
	"skeleton/internal/variable"
	"testing"
	"time"

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

type TransListener struct{}

func (t *TransListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地业务逻辑")
	time.Sleep(1 * time.Second)
	fmt.Println("本地业务逻辑成功")
	// primitive.CommitMessageState 通知 rocketmq 正常提交进topic，不会执行 CheckLocalTransaction
	// primitive.RollbackMessageState 通知 rocketmq 失败，消息丢弃，不会执行 CheckLocalTransaction
	// primitive.UnknowState 通知 rocketmq 异常，让 rocketmq 执行 CheckLocalTransaction
	return primitive.UnknowState
}

func (t *TransListener) CheckLocalTransaction(ext *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("收到Rocketmq主动请求信息,msgID:", ext.MsgId)
	// primitive.CommitMessageState 通知 rocketmq 正常提交进topic
	// primitive.RollbackMessageState 通知 rocketmq 失败，消息丢弃
	return primitive.CommitMessageState
}

func TestMqTransProducer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	rmq, err := variable.MQ.TransProducer(&TransListener{})
	if err != nil {
		t.Errorf("initialization filed:%v", err)
		return
	}
	msg := &primitive.Message{
		Topic: "foo-trans-topic",
		Body:  []byte("test trans message"),
	}
	if err := variable.MQ.SendTransactionMessage(rmq, msg); err != nil {
		t.Errorf("send fail: %s", err)
		return
	}
	t.Log("success")
	<-(chan any)(nil)
}

type FooTransConsumer struct{}

func (*FooTransConsumer) GetTopic() string {
	return "foo-trans-topic"
}

func (*FooTransConsumer) GetSelector() consumer.MessageSelector {
	return consumer.MessageSelector{}
}

func (*FooTransConsumer) Exec() mq.ConsumerExecFunc {
	return func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("message:%v", me)
		return consumer.ConsumeSuccess, nil
	}
}

func TestMqTransConsumer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	go func() {
		if err := variable.MQ.Subscribe(&FooTransConsumer{}); err != nil {
			t.Errorf("fail:%s", err)
		}
	}()
	<-(chan any)(nil)
}
