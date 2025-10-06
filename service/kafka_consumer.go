package service

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(
	brokers *[]string,
	groupID string,
	topic string,
	tlsConfig *tls.Config,
) (*KafkaConsumer, error) {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			CommitInterval: 0, // disable auto-commit
			Brokers:        *brokers,
			GroupID:        groupID, // all consumers with same group id share the messages
			Topic:          topic,
			Dialer: &kafka.Dialer{
				TLS: tlsConfig,
			},
		}),
	}, nil
}

func (c *KafkaConsumer) ListenMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Kafka consumer context done")
			return
		default:
			message, err := c.readMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to read message")
				continue
			}
			log.Info().Msgf(
				"message at offset %d: %s = %s\n",
				message.Offset,
				string(message.Key),
				string(message.Value),
			)
		}
	}
}

func (c *KafkaConsumer) readMessage(ctx context.Context) (*kafka.Message, error) {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return &kafka.Message{}, fmt.Errorf("failed to read message: %w", err)
	}
	err = c.reader.CommitMessages(ctx, message)
	if err != nil {
		return &kafka.Message{}, fmt.Errorf("failed to commit message: %w", err)
	}
	log.Info().Msgf(
		"[x] received message: key=%s, value=%s",
		string(message.Key),
		string(message.Value),
	)
	return &message, nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
