package iokafka

import "github.com/segmentio/kafka-go"

type Message struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
}

type KafkaMessage kafka.Message

func (m Message) ToKafkaMessage() kafka.Message {
	return kafka.Message{
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
		Key:       m.Key,
		Value:     m.Value,
	}
}

func (m KafkaMessage) ToMessage() Message {
	return Message{
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
		Key:       m.Key,
		Value:     m.Value,
	}
}
