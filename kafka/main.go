package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/env"
	// "server/logging"
	"strconv"

	"context"
	"server/models"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	store *models.NotificationStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	log.Info().Msg("Consuming messages from Kafka")

	log.Info().Msgf("Claim Topic: %v", claim.Topic())

	for msg := range claim.Messages() {
		log.Info().Msgf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		userID := string(msg.Key)
		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		log.Info().Msgf("Notification: %v", notification)
		consumer.store.AddNotification(userID, notification)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func FindUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func GetIDFromRequest(r *http.Request, idName string) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{
		fmt.Sprintf("%s:%s", env.DefaultConfig.KAFKA_HOST, env.DefaultConfig.KAFKA_PORT),
	}, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func SetupConsumer() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup([]string{
		fmt.Sprintf("%s:%s", env.DefaultConfig.KAFKA_HOST, env.DefaultConfig.KAFKA_PORT),
	}, env.DefaultConfig.KAFKA_CONSUMER_GROUP, config)
	if err != nil {
		return nil, err
	}

	return consumerGroup, nil
}

func SetupConsumerGroup(store *models.NotificationStore) {
	// func SetupConsumerGroup(ctx context.Context, store *models.NotificationStore) {
	consumerGroup, err := SetupConsumer()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup consumer group")
	}

	log.Info().Msg("Kafka Consumer 🥣 running")

	consumer := &Consumer{store: store}
	ctx := context.Background()

	for {
		err = consumerGroup.Consume(ctx, []string{env.DefaultConfig.KAFKA_TOPIC}, consumer)
		if err != nil {
			log.Error().Err(err).Msg("failed to consume message")
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func SendKafkaMessage(producer sarama.SyncProducer, users []models.User, notification models.Notification) error {

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal notification")
		return errors.New("failed to marshal notification")
	}

	msg := &sarama.ProducerMessage{
		Topic: env.DefaultConfig.KAFKA_TOPIC,
		Key:   sarama.StringEncoder(strconv.Itoa(notification.ToID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
