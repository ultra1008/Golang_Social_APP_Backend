package producer

import (
	"fmt"

	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/internal/cache"
	"github.com/streadway/amqp"
)

type FeedProducer struct {
	ch  *amqp.Channel
	cfg *config.RabbitMQConfig
}

func NewFeedProducer(ch *amqp.Channel, cfg *config.RabbitMQConfig, cache cache.Cache) *FeedProducer {
	return &FeedProducer{
		ch:  ch,
		cfg: cfg,
	}
}

func (f *FeedProducer) SendFeedUpdateMessage(p []byte) error {
	err := f.ch.Publish(
		f.cfg.FeedExchangeName,
		f.cfg.FeedRoutingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         p,
		},
	)
	if err != nil {
		return fmt.Errorf("producer.SendFeedUpdateMessage - can't send message to queue: %v", err)
	}

	return nil
}
