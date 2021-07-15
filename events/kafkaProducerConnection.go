package events

import (
	"github.com/Shopify/sarama"
	"sync"
)

type Connection struct {
	sarama.SyncProducer
}

var kafkaConnection *Connection
var once sync.Once

// GetInstance creates one instance and always returns that one instance
func GetInstance() *Connection {
	// only executes this once
	once.Do(func() {
		_, err := connectProducer()
		if err != nil {
			panic(err)
		}
	})
	return kafkaConnection
}

func connectProducer() (sarama.SyncProducer,error) {
	brokersUrl := []string{"localhost:19092", "localhost:29092", "localhost:39092"}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 7
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		panic(err)
	}

	kafkaConnection = &Connection{conn}

	return kafkaConnection, nil
}
