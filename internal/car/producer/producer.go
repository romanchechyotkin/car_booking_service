package paymentproducer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type PaymentPlacer struct {
	producer        *kafka.Producer
	topic           string
	deliveryChannel chan kafka.Event
}

func NewPaymentPlacer(pr *kafka.Producer, topic string) *PaymentPlacer {
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