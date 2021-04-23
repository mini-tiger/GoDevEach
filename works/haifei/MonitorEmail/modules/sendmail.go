package modules

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"haifei/MonitorEmail/g"
	"mime"
	"time"
)

func SendMail(MailTo string) {
	defer func() {
		if err := recover(); err != nil {
			_ = g.GetLog().Error("%s\n", err)
		}
	}()

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", g.GetConfig().MailUser, "Haier") // 发件人
	msg.SetHeader("To", // 收件人
		msg.FormatAddress(MailTo, "Haier"),
		//m.FormatAddress("********@qq.com", "郭靖"),
	)

	msg.SetHeader("Subject", fmt.Sprintf("运维监控报告-%s", time.Now().Format(g.TimeLayout))) // 主题
	msg.SetBody("text/html", g.HtmlBuffer.String())                                     // 正文

	//添加附件
	name := "运维监控.xlsx"
	msg.Attach(g.Fp,
		gomail.Rename(name),
		gomail.SetHeader(map[string][]string{
			"Content-Disposition": []string{ //中文文件名
				fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", name)),
			},
		}),
	)

	d := gomail.NewDialer(g.GetConfig().MailAddr, g.GetConfig().MailPort, g.GetConfig().MailUser, g.GetConfig().MailPass) // 发送邮件服务器、端口、发件人账号、授权码
	if err := d.DialAndSend(msg); err != nil {
		//panic(err)
		_ = g.GetLog().Warn("%s 邮件发送失败\n", MailTo)
		panic(errors.New(fmt.Sprintf("邮件发送过程失败原因:%s\n", err)))
	} else {
		_ = g.GetLog().Warn("%s 邮件发送成功\n", MailTo)

	}

}
