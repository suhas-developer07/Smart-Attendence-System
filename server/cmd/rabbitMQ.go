package cmd

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitmqConnection struct {
	conn  *amqp.Connection
	chann *amqp.Channel
}

func NewRabbitmqConnection() *RabbitmqConnection {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	log.Print(rabbitmqUrl)
	if rabbitmqUrl == "" {
		log.Fatalln("missing rabbitmq url env variable")
	}

	conn, err := amqp.Dial(rabbitmqUrl)

	if err != nil {
		log.Fatalln("failed to connect rabbitmq broker Error:", err.Error())
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalln("failed to open rabbitmq channel , Error:", err.Error())
	}

	QUEUE_NAME := os.Getenv("QUEUE_NAME")

	_, err = ch.QueueDeclare(
		QUEUE_NAME,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("failed to declare a queue", err)
	}

	log.Println("connected to rabbitMQ")
	return &RabbitmqConnection{
		conn:  conn,
		chann: ch,
	}
}
