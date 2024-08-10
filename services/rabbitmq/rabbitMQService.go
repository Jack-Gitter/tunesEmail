package rabbitmq

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jack-Gitter/tunesEmailService/services/db/user"
	"github.com/Jack-Gitter/tunesEmailService/services/email"
	"github.com/rabbitmq/amqp091-go"
)

type QueueMessageType string 

const (
	POST QueueMessageType = "POST"
)

type RabbitMQPostMessage struct {
    Type QueueMessageType
    Poster string
}

type RabbitMQService struct {
    Conn *amqp091.Connection
    Chan *amqp091.Channel
    QName string
    UserService user.IUserService
    EmailService email.IEmailService
}

type IRabbitMQService interface {
    Connect() 
    Read()
}

func(rmq *RabbitMQService) Read() error {

    msgs, err := rmq.Chan.Consume(
      rmq.QName,
      "",     
      true,   
      false,  
      false, 
      false,  
      nil,   
    )

    if err != nil {
        return err
    }

    var forever chan struct{}

    go func() {
      for d := range msgs {
          postMessage := &RabbitMQPostMessage{}
          json.Unmarshal(d.Body, postMessage)
          emails, err := rmq.UserService.GetUserFollowerEmails(postMessage.Poster)
          if err != nil {
              fmt.Println(err.Error())
          }
          fmt.Println(emails)
          err = rmq.EmailService.SendEmail(emails, []byte("test!"))
          if err != nil {
              fmt.Println(err.Error())
          }
      }
    }()

    fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever

    return nil
}

func(rmq *RabbitMQService) Connect() {
    connectionString := os.Getenv("RABBIT_MQ_CONNECTION_STRING")

    conn, err := amqp091.Dial(connectionString)

    if err != nil {
        panic(err.Error())
    }

    rmq.Conn = conn

    ch, err := conn.Channel()

    if err != nil {
        panic(err.Error())
    }

    rmq.Chan = ch
    rmq.QName = "emailQueue"

    _, err = ch.QueueDeclare(
      rmq.QName, 
      false,   
      false,   
      false,  
      false, 
      nil,  
    )
}

