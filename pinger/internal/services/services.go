package services

import (
	"encoding/json"
	"fmt"
	"log"
	"pinger/internal/config"
	"pinger/internal/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PostData(containersL []models.Container) {
	c := config.New().RabbitMQConfig
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq:%s/", c.RABBITMQ_DEFAULT_USER, c.RABBITMQ_DEFAULT_PASS, c.RABBITMQ_PORT))
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ from pinger: %s;\nTry again", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

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
	jsonData, err := json.Marshal(containersL)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	err = ch.Publish(
		"",     // обменник
		q.Name, // имя очереди
		false,  // обязательное
		false,  // немедленное
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(jsonData),
		})
	if err != nil {
		log.Printf("Failed to publish a message: %s", err)
		return
	}
	log.Printf("Sent message: %s", jsonData)
}
