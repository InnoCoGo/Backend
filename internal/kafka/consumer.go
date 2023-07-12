package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var consumer sarama.Consumer

type ConsumerCallback func(data []byte)

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

func ConsumerMsg(callBack ConsumerCallback) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if nil != err {
		logrus.Errorf("consume partition error: %s", err.Error())
		return
	}

	defer partitionConsumer.Close()
	for {
		msg := <-partitionConsumer.Messages()
		if nil != callBack {
			callBack(msg.Value)
		}
	}
}

func CloseConsumer() {
	if nil != consumer {
		consumer.Close()
	}
}
