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

type RabbitMQMessage struct {
    Type QueueMessageType
}

type RabbitMQPostMessage struct {
    RabbitMQMessage
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

    go rmq.readFunc(msgs)
    <-forever

    return nil
}


func (rmq *RabbitMQService) readFunc(msgs <-chan amqp091.Delivery) {
      for d := range msgs {
          t, err := rmq.getMessageType(d)
          if err != nil {
             continue
          }
          switch t {
            case POST:
              postMessage := &RabbitMQPostMessage{}
              json.Unmarshal(d.Body, postMessage)
              rmq.handlePostMessage(postMessage)
          }
      }
}

func(rmq *RabbitMQService) getMessageType(d amqp091.Delivery) (QueueMessageType, error) {
    message := &RabbitMQMessage{}
    err := json.Unmarshal(d.Body, message)
    if err != nil {
        return "", err
    }
    return message.Type, nil
}

func (rmq *RabbitMQService) handlePostMessage(postMessage *RabbitMQPostMessage) {
      emails, err := rmq.UserService.GetUserFollowerEmails(postMessage.Poster)
      if err != nil {
          panic(err.Error())
      }
      username, err := rmq.UserService.GetUsername(postMessage.Poster)
      if err != nil {
          panic(err.Error())
      }
      msg := []byte(fmt.Sprintf("%s Has posted a new post! go check it out", username))
      err = rmq.EmailService.SendEmail(emails, msg)
      if err != nil {
          panic(err.Error())
      }
}
