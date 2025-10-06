package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"my_rest_server/security"
	"my_rest_server/service"

	"github.com/rs/zerolog/log"
)

const (
	userTopicGroup   = KafkaTopicUsers + "-group"
	closableCapacity = 4
)

type AppConfig struct {
	kafkaTlsConfig *tls.Config
	closables      []Closable
}

func SetupApp() *AppConfig {
	ctx := context.Background()
	closables := make([]Closable, 0, closableCapacity)

	kafkaTlsConfig, err := security.NewtlsConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load TLS config: %v", err))
	}
	// closables = append(closables, NewKafkaProducer(kafkaTlsConfig, KafkaTopicUsers))
	closables = append(closables, newKafkaConsumer(ctx, kafkaTlsConfig, "users", userTopicGroup))

	return &AppConfig{
		kafkaTlsConfig: kafkaTlsConfig,
		closables:      closables,
	}
}

func newKafkaConsumer(
	ctx context.Context,
	kafkaTlsConfig *tls.Config,
	topic string,
	groupId string,
) Closable {
	service, err := service.NewKafkaConsumer(
		&[]string{KafkaBrokerAddress},
		groupId,
		topic,
		kafkaTlsConfig,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create Kafka consumer: %v", err))
	}
	go service.ListenMessages(ctx)
	return service
}

func (a *AppConfig) NewKafkaProducer(topic string) *service.KafkaProducer {
	kafkaProducer, err := service.NewKafkaProducer(
		KafkaBrokerAddress,
		topic,
		a.kafkaTlsConfig,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create Kafka producer: %v", err))
	}
	a.closables = append(a.closables, kafkaProducer)

	return kafkaProducer
}

func (a *AppConfig) Close() {
	for _, closable := range a.closables {
		err := closable.Close()
		if err != nil {
			log.Error().Err(err).Msgf("failed to close %T", closable)
		}
	}
}

type Closable interface {
	Close() error
}
