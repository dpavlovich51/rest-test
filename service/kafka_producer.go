package service

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker string, topic string, tlsConfig *tls.Config) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:      kafka.TCP(broker),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		Transport: &kafka.Transport{TLS: tlsConfig},
	}
	return &KafkaProducer{writer: writer}, nil
}

func (p *KafkaProducer) ProduceMessage(value string) error {
	return p.ProduceMessageWithKey("", value)
}

// key - message key. it effects on partitioning. messages with same key go to same partition
// value - message value
func (p *KafkaProducer) ProduceMessageWithKey(key string, value string) error {
	var msg kafka.Message
	if key == "" {
		msg = kafka.Message{Value: []byte(value)}
	} else {
		msg = kafka.Message{Key: []byte(key), Value: []byte(value)}
	}
	log.Info().Msgf("[x] send message: key=%s, value=%s", key, value)

	err := p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
