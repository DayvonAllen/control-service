package event_consumer

import (
	"context"
	"example.com/app/domain"
	"example.com/app/repo"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

func KafkaConsumerGroup() {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	group := "go-kafka-consumer"
	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}

	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())

	client, err := sarama.NewConsumerGroup(brokers, group, config)

	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	topics := []string{"user"}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		user := new(domain.Message)
		err := msgpack.Unmarshal(message.Value, user)
		log.Printf("Message claimed: value = %v, timestamp = %v, topic = %s", user, message.Timestamp, message.Topic)

		if err != nil {
			return err
		}

		err = repo.ProcessMessage(*user)

		if err != nil {
			return err
		}

		session.MarkMessage(message, "")
	}

	return nil
}
