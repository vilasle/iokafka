package iokafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Partition int
	Topic string
	AttemtsOnFail int
}

type Scanner struct {
	reader *kafka.Reader
	ctx    context.Context
	msg    Message
	err    error
	lastOffset int64
	timeout int
	attemptsOnFail int
}

func NewScanner(config KafkaConfig) Scanner {
	scanner := Scanner{
		ctx:    context.Background(),
		attemptsOnFail: config.AttemtsOnFail,
	}

	scanner.initKafkaReader(config)
	return scanner
}

func (r *Scanner) Scan() bool {
	var err error
	for t := 0; t < r.attemptsOnFail; t++ {
		msg, err := r.reader.ReadMessage(r.ctx)
		if err == nil {
			r.msg = KafkaMessage(msg).ToMessage()
			r.lastOffset = msg.Offset
			return true
		}
		time.Sleep(time.Duration(r.timeout))
	}
	r.err = err
	return false
}

func (r *Scanner) Message() Message {
	return r.msg
}

func (r *Scanner) Err() error {
	return r.err
}

func (r *Scanner) initKafkaReader(config KafkaConfig) {
	
	conf := kafka.ReaderConfig{
		Brokers:   config.Brokers,
		GroupID:   config.GroupID,
		Partition: config.Partition,
		Topic:     config.Topic,
	}

	kafka.NewReader(conf)
}

