package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	"github.com/itoqsky/InnoCoTravel-backend/pkg/protocol"
	"github.com/sirupsen/logrus"
)

var consumer sarama.Consumer

type ConsumerCallback func(msg *protocol.Message)

func InitConsumer(hosts string) {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if nil != err {
		logrus.Errorf("init kafka client error: %s", err.Error())
	}

	consumer, err = sarama.NewConsumerFromClient(client)
	if nil != err {
		logrus.Errorf("init kafka consumer error: %s", err.Error())
	}
}

func Consume(callBack ConsumerCallback) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if nil != err {
		logrus.Errorf("consume partition error: %s", err.Error())
		return
	}

	defer partitionConsumer.Close()
	for {
		rawMsg := <-partitionConsumer.Messages()
		if nil != callBack {
			msg := &protocol.Message{}
			err := proto.Unmarshal(rawMsg.Value, msg)
			if nil != err {
				logrus.Errorf("unmarshal message error: %s", err.Error())
				continue
			}

			callBack(msg)
		}
	}
}

func CloseConsumer() {
	if nil != consumer {
		consumer.Close()
	}
}
