package rabbit

import (
	"github.com/BogdanStaziyev/softcery-test/pkg"
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
			go func() {
				log.Println(mes)
				err = pkg.MakeVariants(mes)
				if err != nil {
					log.Println(err)
					return
				}
			}()
		}
	}()
	<-forever
	return nil
}
