package emailproducer

import "fmt"

func ignoring() {
	fmt.Println("ignoring")
}

//
//import (
//	"fmt"
//	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//	"log"
//)
//
//type EmailPlacer struct {
//	producer        *kafka.Producer
//	topic           string
//	deliveryChannel chan kafka.Event
//}
//
//func NewEmailPlacer(p *kafka.Producer, topic string) *EmailPlacer {
//	return &EmailPlacer{
//		producer:        p,
//		topic:           topic,
//		deliveryChannel: make(chan kafka.Event),
//	}
//}
//
//func (ep *EmailPlacer) SendEmail(email, emailType string) error {
//	format := fmt.Sprintf("%s %s", email, emailType)
//	payload := []byte(format)
//
//	msg := &kafka.Message{
//		TopicPartition: kafka.TopicPartition{Topic: &ep.topic, Partition: kafka.PartitionAny},
//		Value:          payload,
//	}
//
//	err := ep.producer.Produce(msg, ep.deliveryChannel)
//	if err != nil {
//		return err
//	}
//
//	<-ep.deliveryChannel
//	log.Printf("placed payload to kafka queue: %s\n", payload)
//	return nil
//}
