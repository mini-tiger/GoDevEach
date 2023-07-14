package main

/**
 * @Author: Tao Jun
 * @Since: 2023/7/13
 * @Desc: 163demo.go
**/

import (
	"crypto/tls"
	gomail "gopkg.in/gomail.v2"
)

func main() {

	msg := gomail.NewMessage()
	msg.SetHeader("From", "61566027@163.com")
	msg.SetHeader("To", "61566027@163.com")
	msg.SetHeader("Subject", "测试邮件")
	msg.SetBody("text/html", "<b>This is the body of the mail</b>")
	// msg.Attach("/home/User/cat.jpg")

	//n := gomail.NewDialer("smtp.163.com", 465, "61566027", "ETPOQJQYTBJLQWDT")
	//n := gomail.NewDialer("172.22.16.58", 465, "61566027", "ETPOQJQYTBJLQWDT")
	n := gomail.NewDialer("smtp.163.com", 25, "61566027", "ETPOQJQYTBJLQWDT")
	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
}
