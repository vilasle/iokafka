package iokafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type ScannerConfig struct {
	Brokers       []string
	GroupID       string
	Partition     int
	Topic         string
	AttemtsOnFail int
	FailTimeout   int
}

type Scanner struct {
	reader     *kafka.Reader
	ctx        context.Context
	msg        Message
	err        error
	lastOffset int64
	cnf        ScannerConfig
}

func NewScanner(config ScannerConfig) Scanner {

	s := Scanner{
		ctx: context.Background(),
		cnf: config,
	}

	if s.cnf.AttemtsOnFail == 0 {
		s.cnf.AttemtsOnFail = 1
	}

	if s.cnf.FailTimeout == 0 {
		s.cnf.FailTimeout = 10
	}

	s.initKafkaReader()
	return s
}

func (s *Scanner) Scan() bool {
	var err error
	for t := 0; t < s.cnf.AttemtsOnFail; t++ {
		msg, err := s.reader.ReadMessage(s.ctx)
		if err == nil {
			s.msg = KafkaMessage(msg).ToMessage()
			s.lastOffset = msg.Offset
			return true
		}
		time.Sleep(time.Duration(s.cnf.FailTimeout))
	}
	s.err = err
	return false
}

func (s *Scanner) Message() Message {
	return s.msg
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) initKafkaReader() {

	conf := kafka.ReaderConfig{
		Brokers:   s.cnf.Brokers,
		GroupID:   s.cnf.GroupID,
		Partition: s.cnf.Partition,
		Topic:     s.cnf.Topic,
	}

	s.reader = kafka.NewReader(conf)
}
