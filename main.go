package main

import (
	"github.com/Jack-Gitter/tunesEmailService/services/db/general"
	"github.com/Jack-Gitter/tunesEmailService/services/db/user"
	"github.com/Jack-Gitter/tunesEmailService/services/email"
	"github.com/Jack-Gitter/tunesEmailService/services/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        panic("bad")
    }

    emailService := &email.EmailService{}
    emailService.Authenticate()
    conn := general.ConnectToDB()
    userDao := &user.UserDAO{}
    userService := user.UserService{DB: conn, UserDAO: userDao}
    rabbitMQService := rabbitmq.RabbitMQService{}
    rabbitMQService.UserService = &userService
    rabbitMQService.EmailService = emailService
    rabbitMQService.Connect()
    rabbitMQService.Read()


}
