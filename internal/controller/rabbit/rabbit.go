package rabbit

import (
	// rabbit
	"github.com/streadway/amqp"

	// external
	"github.com/BogdanStaziyev/softcery-test/pkg/logger"
	"github.com/BogdanStaziyev/softcery-test/pkg/utils"
)

const queueName = "image"

type Rabbit struct {
	conn *amqp.Connection
	l    logger.Interface
}

func NewRabbit(url string, l logger.Interface) Rabbit {
	conn, err := amqp.Dial(url)
	if err != nil {
		l.Fatal("Unable to create new RabbitMQ connection: ", err)
	}
	return Rabbit{
		conn: conn,
	}
}

// CreateQueue creating a queue using a constant name
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

// PublishImage writes the file path to the queue
func (r *Rabbit) PublishImage(path string) error {
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
func (r *Rabbit) Consumer() error {
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
			err = utils.MakeVariants(mes)
			if err != nil {
				r.l.Error(err, "- rabbit")
				return
			}
		}
	}()
	<-forever
	return nil
}
