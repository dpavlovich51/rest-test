package config

import (
	"context"
	"crypto/tls"
	"fmt"
	c "my_rest_server/client"
	"my_rest_server/security"
	s "my_rest_server/service"

	// "my_rest_server/service"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	userTopicGroup   = KafkaTopicUsers + "-group"
	closableCapacity = 4
)

type AppConfig struct {
	RouterHandler http.Handler

	kafkaTlsConfig *tls.Config
	closables      []Closable
}

func SetupApp() *AppConfig {
	// Initialize logger
	SetupLogger()
	// Initialize closables slice
	closables := make([]Closable, 0, closableCapacity)
	// Initialize context
	ctx := context.Background()
	// Initialize TLS config for Kafka
	kafkaTlsConfig, err := security.NewTLSConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load TLS config: %v", err))
	}
	// Initialize UserService
	userService, err := s.NewUserService(
		RedisAddress,
		RedisPassword,
		RedisDB,
		createKafkaProducer(kafkaTlsConfig), // pass new kafka producer to user service
	)
	if err != nil {
		panic(fmt.Errorf("failed to create UserService: %v", err))
	} else {
		closables = append(closables, userService)
	}
	// Initialize routerHandler
	routerHandler := WrapWithLogging(SetupRouter(userService))

	// closables = append(closables, NewKafkaProducer(kafkaTlsConfig, KafkaTopicUsers))
	closables = append(closables, newKafkaConsumer(ctx, kafkaTlsConfig, "users", userTopicGroup))

	return &AppConfig{
		RouterHandler:  routerHandler,
		kafkaTlsConfig: kafkaTlsConfig,
		closables:      closables,
	}
}

func createKafkaProducer(kafkaTlsConfig *tls.Config) *c.KafkaProducer {
	producer, err := c.NewKafkaProducer(KafkaBrokerAddress, KafkaTopicUsers, kafkaTlsConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create Kafka producer: %v", err))
	}
	return producer
}

func newKafkaConsumer(
	ctx context.Context,
	kafkaTlsConfig *tls.Config,
	topic string,
	groupId string,
) Closable {
	service, err := c.NewKafkaConsumer(
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

func newKafkaProducer(
	kafkaTlsConfig *tls.Config,
	topic string,
	closables *[]Closable,
) *c.KafkaProducer {
	kafkaProducer, err := c.NewKafkaProducer(
		KafkaBrokerAddress,
		topic,
		kafkaTlsConfig,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create Kafka producer: %v", err))
	}
	*closables = append(*closables, kafkaProducer)
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
