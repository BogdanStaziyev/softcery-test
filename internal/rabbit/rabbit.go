package rabbit

import (
	"github.com/streadway/amqp"
	"log"
)

const queueName = "image"

type Rabbit struct {
	conn *amqp.Connection
}

func NewRabbit(url string) Rabbit {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Unable to create new RabbitMQ connection: %q\n", err)
	}
	return Rabbit{
		conn: conn,
	}
}

func (r *Rabbit) CreateQueue() error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rabbit) PublishImage(name string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(name),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Rabbit) Consumer(path string, imageName string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	message, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)

	go func() {
		for data := range message {
			mes := string(data.Body)
			go func() {
				log.Println(mes)
				return
			}()
		}
	}()
	<-forever
	return nil
}
