package iokafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	writer *kafka.Writer
}

func NewWriter(brokerUrl string, topic string) *Writer {
	w := Writer{}
	w.initKafkaWriter(brokerUrl, topic)

	return &w
}

func (w *Writer) Write(msg Message) (err error) {
	return w.writer.WriteMessages(context.Background(), msg.ToKafkaMessage())
}

func (w *Writer) Close() (err error) {
	return err
}

func (w Writer) initKafkaWriter(kafkaURL string, topic string) {
	w.writer = &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
