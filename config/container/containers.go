package container

import (
	"github.com/BogdanStaziyev/softcery-test/config"
	"github.com/BogdanStaziyev/softcery-test/internal/app"
	"github.com/BogdanStaziyev/softcery-test/internal/infra/handlers"
	"github.com/BogdanStaziyev/softcery-test/internal/rabbit"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
)

type Container struct {
	Services
	Controllers
	Queue
}

type Services struct {
	app.ImageService
}

type Controllers struct {
	handlers.ImageHandler
}

type Queue struct {
	rabbit.Rabbit
}

func New(conf config.Configuration) Container {
	_ = getDbSess(conf)

	rabbitMQ := rabbit.NewRabbit(conf.RabbitURL)
	imageService := app.NewImageService(conf.FileStorageLocation)
	imageHandler := handlers.NewImageHandler(imageService)

	//err := rabbitMQ.CreateQueue()
	//if err != nil {
	//	log.Println(err)
	//}
	//err = rabbitMQ.PublishImage("hallo, halo")
	//if err != nil {
	//	log.Println(err)
	//}
	//time.Sleep(time.Second * 2)
	//err = rabbitMQ.Consumer(conf.FileStorageLocation, "name.png")
	//if err != nil {
	//	log.Println(err)
	//}
	return Container{
		Services: Services{
			imageService,
		},
		Controllers: Controllers{
			imageHandler,
		},
		Queue: Queue{
			rabbitMQ,
		},
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

//func getRabbit(conf config.Configuration) *amqp.Connection {
//	conn, err := amqp.Dial(conf.RabbitURL)
//	if err != nil {
//		log.Fatalf("Unable to create new RabbitMQ connection: %q\n", err)
//	}
//	ch, err := conn.Channel()
//	if err != nil {
//		log.Fatalf("Unable to create connection sever channel RabbitMQ: %q\n", err)
//	}
//	q, err := ch.QueueDeclare("testQueue", false, false, false, false, nil)
//	if err != nil {
//		log.Fatalf("Unable to declares a queue RabbitMQ: %q\n", err)
//	}
//	fmt.Println(q)
//	err = ch.Publish("", "testQueue", false, false, amqp.Publishing{
//		ContentType: "text/plain",
//		Body:        []byte("hello"),
//	})
//	if err != nil {
//		fmt.Println(err)
//	}
//	return conn
//}
//
//func getConsumer(conf config.Configuration) {
//	conn, err := amqp.Dial(conf.RabbitURL)
//	if err != nil {
//		log.Fatalf("Unable to create new consumer RabbitMQ connection: %q\n", err)
//	}
//	ch, err := conn.Channel()
//	if err != nil {
//		log.Fatalf("Unable to create connection sever channel consumer RabbitMQ: %q\n", err)
//	}
//	msg, err := ch.Consume("testQueue", "", true, false, false, false, nil)
//	reader := make(chan bool)
//	go func() {
//		for d := range msg {
//			fmt.Printf("message: %s\n", d.Body)
//		}
//	}()
//	<-reader
//}
