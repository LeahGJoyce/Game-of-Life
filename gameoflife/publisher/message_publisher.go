package publisher

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

var env *stream.Environment

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func GetEnv() *stream.Environment {
	if env == nil {
		host := os.Getenv("RABBITMQ_HOST")
		if host == "" {
			host = "localhost" // default host if not set
		}
		portStr := os.Getenv("RABBITMQ_PORT")
		if portStr == "" {
			portStr = "5552" // default port if not set
		}
		port, err := strconv.Atoi(portStr) // Convert port from string to int
		CheckErr(err)

		user := os.Getenv("RABBITMQ_USER")
		if user == "" {
			user = "user" // default user if not set
		}
		password := os.Getenv("RABBITMQ_PASS")
		if password == "" {
			password = "password" // default password if not set
		}

		newEnv, err := stream.NewEnvironment(
			stream.NewEnvironmentOptions().
				SetHost(host).
				SetPort(port).
				SetUser(user).
				SetPassword(password).
				SetMaxProducersPerClient(10))
		CheckErr(err)
		env = newEnv
	}
	return env
}

func PublishMessage(env *stream.Environment, streamName, message string) {
	err := env.DeclareStream(streamName, &stream.StreamOptions{MaxLengthBytes: stream.ByteCapacity{}.GB(2)})
	if err != nil {
		fmt.Printf("Error declaring stream '%s': %v\n", streamName, err)
	}

	producer, err := env.NewProducer(streamName, nil)
	if err != nil {
		fmt.Printf("Error creating producer for stream '%s': %v\n", streamName, err)
	}

	err = producer.Send(amqp.NewMessage([]byte(message)))
	if err != nil {
		fmt.Printf("Error sending message to stream '%s': %v\n", streamName, err)
	}

	err = producer.Close()
	if err != nil {
		fmt.Printf("Error closing producer for stream '%s': %v\n", streamName, err)
	}
}
