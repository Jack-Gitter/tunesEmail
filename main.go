package main

import (
	"github.com/Jack-Gitter/tunesEmailService/services/db/general"
	"github.com/Jack-Gitter/tunesEmailService/services/db/user"
	"github.com/Jack-Gitter/tunesEmailService/services/rabbitmq"
)

func main() {


    conn := general.ConnectToDB()
    userDao := &user.UserDAO{}
    userService := user.UserService{DB: conn, UserDAO: userDao}
    rabbitMQService := rabbitmq.RabbitMQService{}
    rabbitMQService.Connect()
    rabbitMQService.UserService = &userService
    rabbitMQService.Read()

}
