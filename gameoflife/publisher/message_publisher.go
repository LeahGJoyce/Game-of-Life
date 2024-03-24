package publisher

import (
	"fmt"
	"log"
	"sync"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
)

var (
	env       *stream.Environment
	producers = make(map[string]*stream.Producer)
	mu        sync.Mutex // To synchronize access to the producers map
)

func setupEnvironment() {
	var err error
	if env == nil {
		env, err = stream.NewEnvironment(
			stream.NewEnvironmentOptions().SetUri("amqp://user:password@rabbitmq:5672/"),
		)
		if err != nil {
			log.Fatalf("Failed to create an environment: %s", err)
		}
	}
}

func getProducer(gameId string) *stream.Producer {
	mu.Lock()
	defer mu.Unlock()

	producer, exists := producers[gameId]
	if !exists {
		var err error
		// Ensure the stream exists for the gameId
		err = env.DeclareStream(gameId, &stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		})
		if err != nil {
			log.Fatalf("Failed to declare a stream: %s", err)
		}

		producer, err = env.NewProducer(gameId, nil)
		if err != nil {
			log.Fatalf("Failed to create a producer: %s", err)
		}
		producers[gameId] = producer
	}

	return producer
}

func PublishMessage(gameId, message string) {
	setupEnvironment()
	producer := getProducer(gameId)

	// Publish a message to the stream
	err := producer.Send(amqp.NewMessage([]byte(message)))
	if err != nil {
		log.Fatalf("Failed to send a message: %s", err)
	}

	fmt.Printf("Published message to stream %s: %s\n", gameId, message)
}


// CleanupProducers closes all producers. Should be called before shutting down the application.
func CleanupProducers() {
	mu.Lock()
	defer mu.Unlock()

	for gameID, producer := range producers {
		err := producer.Close()
		if err != nil {
			log.Printf("Failed to close producer for game %s: %s", gameID, err)
		} else {
			log.Printf("Closed producer for game %s", gameID)
		}
		// Delete the producer from the map to prevent further use.
		delete(producers, gameID)
	}
}

// CleanupEnvironment closes the RabbitMQ stream environment. Should be called after all producers are closed.
func CleanupEnvironment() {
	if env != nil {
		err := env.Close()
		if err != nil {
			log.Fatalf("Failed to close stream environment: %s", err)
		}
		env = nil // Ensure the reference is removed to prevent reuse.
		log.Println("Closed stream environment")
	}
}
