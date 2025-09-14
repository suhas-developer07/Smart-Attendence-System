package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitmqRepo struct {
	conn  *amqp.Connection
	chann *amqp.Channel
}

func NewRabbitmqRepo(conn *amqp.Connection, chann *amqp.Channel) *RabbitmqRepo {
	return &RabbitmqRepo{
		conn:  conn,
		chann: chann,
	}
}

type Data struct {
	USN    string   `json:"usn"`
	Images []string `json:"images"`
}

func (rmq *RabbitmqRepo) SendMessage(data Data) error {
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		return fmt.Errorf("missing QUEUE_NAME env variable")
	}

	_, err := rmq.chann.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	err = rmq.chann.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Println("Messege sent Successfully.")

	return nil
}
