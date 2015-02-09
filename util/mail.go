package util

import (
    "gopkg.in/gomail.v1"
)

func SendMail(data string) {
    msg := gomail.NewMessage()
    msg.SetHeader("From", "zhangtian0809@gmail.com")
    msg.SetHeader("To", "zhangtian0809@gmail.com")
    msg.SetHeader("Subject", "Market Data")
    msg.SetBody("text/html", data)

    mailer := gomail.NewMailer("smtp.gmail.com", "fatbirdstock@gmail.com", "lovelybird", 465)
    if err := mailer.Send(msg); err != nil {
        panic(err)
    }
}
