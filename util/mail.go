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

    //f, err := gomail.OpenFile("/home/Alex/lolcat.jpg")
    //if err != nil {
     //   panic(err)
    //}
    //msg.Attach(f)

    // Send the email to Bob, Cora and Dan
    mailer := gomail.NewMailer("smtp.gmail.com", "fatbirdstock@gmail.com", "lovelybird", 465)
    if err := mailer.Send(msg); err != nil {
        panic(err)
    }
}