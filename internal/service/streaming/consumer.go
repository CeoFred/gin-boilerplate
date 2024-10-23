package streaming

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type EventConsumers map[string]func(msg *sarama.ConsumerMessage) error

type Consumer struct {
	ready          chan bool
	client         *sarama.ConsumerGroup
	ctx            context.Context
	isPaused       *bool
	EventConsumers EventConsumers
}

func NewConsumer(cfg *Config) (*Consumer, error) {
	// Create consumer instance
	consumer := &Consumer{
		ready:          make(chan bool),
		ctx:            cfg.Ctx,
		EventConsumers: make(EventConsumers),
	}

	// Setup and return any errors
	if err := consumer.setUp(cfg); err != nil {
		return nil, err
	}
	return consumer, nil
}

func (c *Consumer) setUp(cfg *Config) error {
	log.Println("Starting a new Sarama consumer")
	if cfg.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return fmt.Errorf("error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version

	switch cfg.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		return fmt.Errorf("unrecognized consumer group partition assignor: %s", cfg.Assignor)
	}

	if cfg.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	client, err := sarama.NewConsumerGroup(strings.Split(cfg.Brokers, ","), cfg.Group, config)
	if err != nil {
		return fmt.Errorf("error creating consumer group client: %v", err)
	}

	c.client = &client
	return nil
}

func (c *Consumer) ToggleConsumptionFlow() {
	client := *c.client
	if *c.isPaused {

		client.PauseAll()
		*c.isPaused = false
		log.Println("Resuming consumption")
	} else {

		client.ResumeAll()
		*c.isPaused = true
		log.Println("Pausing consumption")
	}
}

func (consumer *Consumer) Consume(topics string, messageHandler func(msg *sarama.ConsumerMessage) error) error {

	for _, topic := range strings.Split(topics, ",") {
		consumer.EventConsumers[topic] = messageHandler
	}

	client := *consumer.client

	if consumer.client == nil {
		return errors.New("consumer client is not initialized")
	}

	consumer.ready = make(chan bool)

	go func() {
		for {
			if err := client.Consume(consumer.ctx, strings.Split(topics, ","), consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Printf("Consumer group has been closed")
					return
				}
				log.Printf("Error from consumer: %v", err)
				return
			}

			if consumer.ctx.Err() != nil {
				log.Printf("Context cancelled: %v", consumer.ctx.Err())
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	return nil
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel was closed")
				return nil
			}
			if err := consumer.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
			}
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (consumer *Consumer) processMessage(msg *sarama.ConsumerMessage) error {

	if consumer.EventConsumers[msg.Topic] != nil {
		go consumer.EventConsumers[msg.Topic](msg)
	}
	return nil
}
