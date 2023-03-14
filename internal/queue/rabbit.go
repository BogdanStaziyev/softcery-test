package queue

import (
	// queue
	"github.com/streadway/amqp"

	// external
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"

	// internal
	"github.com/BogdanStaziyev/softcery-test/internal/usecase"
)

const queueName = "image"

type rabbit struct {
	conn *amqp.Connection
	str  usecase.Storage
	l    logger.Interface
}

func NewRabbit(url string, l logger.Interface, storage usecase.Storage) *rabbit {
	conn, err := amqp.Dial(url)
	if err != nil {
		l.Fatal("Unable to create new RabbitMQ connection: ", err)
	}
	return &rabbit{
		conn: conn,
		l:    l,
		str:  storage,
	}
}

// CreateQueue creating a queue using a constant name
func (r *rabbit) CreateQueue() error {
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

// PublishImage writes the file path to the queue
func (r *rabbit) PublishImage(path string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	//Send current image path to rabbitMQ
	err = ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(path),
	})
	if err != nil {
		return err
	}
	return nil
}

// Consumer read the file path and sends the file to create versions with a smaller size
func (r *rabbit) Consumer() error {
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

	//Read image path and create different quality
	go func() {
		for data := range message {
			mes := string(data.Body)
			//Make variants 75%, 50%, 25% size image
			err = r.str.MakeVariants(mes)
			if err != nil {
				r.l.Error(err, "- queue")
				return
			}
		}
	}()
	<-forever
	return nil
}
