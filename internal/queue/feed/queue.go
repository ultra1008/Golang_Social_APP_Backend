package feed

import (
	"fmt"
	"log"

	"github.com/niklod/highload-social-network/config"
	"github.com/streadway/amqp"
)

func NewQueueChannel(conn *amqp.Connection, cfg *config.RabbitMQConfig) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("feed.NewQueue - can't get rabbitmq channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		cfg.FeedExchangeName, // exchange name
		"direct",             // exchange kind
		true,                 // durable
		false,                // auto delete
		false,                // internal
		false,                // no wait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("feed.NewQueue - can't create exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		cfg.FeedQueueName, // queue name
		true,              // durable
		false,             // auto delete
		false,             // exclusive
		false,             // no wait
		nil,
	)

	err = ch.QueueBind(cfg.FeedQueueName, cfg.FeedRoutingKey, cfg.FeedExchangeName, false, nil)
	if err != nil {
		return nil, fmt.Errorf("feed.NewQueue - can't bind queue to exchange: %v", err)
	}

	log.Printf("queue %q binded to %q exchange\n", cfg.FeedQueueName, cfg.FeedExchangeName)

	return ch, nil
}
