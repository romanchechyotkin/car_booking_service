package paymentproducer

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ignoring() {
	fmt.Print("ignoring")
}

type PaymentPlacer struct {
	producer        *kafka.Producer
	topic           string
	deliveryChannel chan kafka.Event
}

type PaymentPlacerer interface {
	SendPayment(payload []byte) error
}

func NewPaymentPlacer(pr *kafka.Producer, topic string) PaymentPlacerer {
	return &PaymentPlacer{
		producer:        pr,
		topic:           topic,
		deliveryChannel: make(chan kafka.Event),
	}
}

func (p *PaymentPlacer) SendPayment(payload []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}

	err := p.producer.Produce(msg, p.deliveryChannel)
	if err != nil {
		log.Printf("kafka error %v", err)
		return err
	}

	<-p.deliveryChannel
	log.Printf("placed payload to kafka queue: %s\n", payload)
	return nil
}
