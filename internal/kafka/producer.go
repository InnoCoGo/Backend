package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var producer sarama.AsyncProducer
var topic string = "default_message"

func InitProducer(topicInput, hosts string) {
	topic = topicInput
	config := sarama.NewConfig()
	config.Producer.Compression = sarama.CompressionGZIP
	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if nil != err {
		logrus.Errorf("init kafka client error: %s", err.Error())
	}

	producer, err = sarama.NewAsyncProducerFromClient(client)
	if nil != err {
		logrus.Errorf("init kafka async client error: %s", err.Error())
	}
}

func Produce(data []byte) {
	be := sarama.ByteEncoder(data)
	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: be}
}

func Close() {
	if producer != nil {
		producer.Close()
	}
}
