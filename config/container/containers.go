package container

import (
	"fmt"
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/streadway/amqp"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
)

type Container struct {
	Services
	Controllers
}

type Services struct {
}

type Controllers struct {
}

func New(conf config.Configuration) Container {
	_ = getDbSess(conf)
	_ = getRabbit(conf)

	return Container{
		Services:    Services{},
		Controllers: Controllers{},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}

func getRabbit(conf config.Configuration) *amqp.Connection {
	conn, err := amqp.Dial(conf.RabbitURL)
	if err != nil {
		log.Fatalf("Unable to create new RabbitMQ connection: %q\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Unable to create connection sever channel RabbitMQ: %q\n", err)
	}
	q, err := ch.QueueDeclare("testQueue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Unable to declares a queue RabbitMQ: %q\n", err)
	}
	fmt.Println(q)
	err = ch.Publish("", "testQueue", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello"),
	})
	if err != nil {
		fmt.Println(err)
	}
	return conn
}
