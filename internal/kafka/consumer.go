package kafka

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/sirupsen/logrus"
)

var consumer sarama.Consumer

type ConsumerCallback func(msg *core.Message)

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
			msg := &core.Message{}
			// proto.Unmarshal(rawMsg.Value, msg)
			err = json.Unmarshal(rawMsg.Value, msg)
			if nil != err {
				logrus.Errorf("unmarshal message error: %s", err.Error())
				continue
			}
			fmt.Printf("CONSUMERING: %v", msg)

			callBack(msg)
		}
	}
}

func CloseConsumer() {
	if nil != consumer {
		consumer.Close()
	}
}
