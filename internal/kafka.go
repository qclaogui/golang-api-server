package internal

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// KafkaConfig  contains information
type KafkaConfig struct {
	Hosts string
	Topic string
}

// Write2Kafka is used for sending logs to kafka.
type Write2Kafka struct {
	config   *KafkaConfig
	producer sarama.AsyncProducer
}

func (o *Write2Kafka) Write(bt []byte) (n int, err error) {
	o.producer.Input() <- &sarama.ProducerMessage{
		Topic: o.config.Topic,
		Value: sarama.ByteEncoder(bt),
	}

	return len(bt), nil
}

func NewLog2Kafka(cfg *KafkaConfig) io.Writer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 300 * time.Millisecond
	config.Producer.Retry.Max = 10
	//config.Producer.Return.Successes = true

	var brokerList []string
	for _, broker := range strings.Split(cfg.Hosts, ",") {
		if strings.Index(broker, ":") == -1 {
			broker += ":9092"
		}
		brokerList = append(brokerList, broker)
	}

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Kafka producer:", err)
	}

	kaf := &Write2Kafka{cfg, producer}

	go func() {
		for err := range kaf.producer.Errors() {
			fmt.Println("Failed to write access log entry:", err)
		}
	}()

	// only when config.Producer.Return.Successes = true
	//go func() {
	//	for res := range kaf.producer.Successes() {
	//		fmt.Printf("write: %#v\n", res)
	//	}
	//}()

	return kaf
}
