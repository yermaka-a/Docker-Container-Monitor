package services

import (
	"back/internal/config"
	"back/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnContainersRMQ() (*amqp.Connection, *amqp.Channel) {
	c := config.New().RabbitMQConfig
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error
	go func() {
		for {
			conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq:%s/", c.RABBITMQ_DEFAULT_USER, c.RABBITMQ_DEFAULT_PASS, c.RABBITMQ_PORT))
			if err != nil {
				log.Printf("Failed to connect to RabbitMQ: %s;\nTry again", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		ch, err = conn.Channel()
		if err != nil {
			log.Panicf("Failed to open a channel: %s", err)
		}

		q, err := ch.QueueDeclare(
			"containers_queue",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Panicf("Failed to declare a queue: %s", err)
		}

		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Panicf("Failed to register a consumer: %s", err)
		}
		go func() {
			var containers []db.Container
			for msg := range msgs {
				err := json.Unmarshal(msg.Body, &containers)
				if err != nil {
					log.Panicf("Failed unmarshal json: %v", err)
				}
				log.Printf("Received a message: %s", msg.Body)
				db.UpdateContainers(&containers)
			}
		}()
	}()

	return conn, ch
}
