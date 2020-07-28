package mailers

import (
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
)
// DefaultSend 默认未登录时发送
func DefaultSend(name, to, content string) error {
	m := mail.NewMessage()
	m.Subject = fmt.Sprintf("来自: %s", name)
	m.From = envy.Get("SMTP_USER", "")
	m.To = []string{to}
	data := render.Data{
		"content": content,
		"name":    name,
		"to":      to,
	}
	err := m.AddBody(r.HTML("send.html"), data)
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
// SendByLogin 登录后发送邮件
func SendByLogin(name, passwd, to, content string) error {
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := name
	password := passwd
	var err error
	fmt.Println(port, host, user, password)
	smtp, err = mail.NewSMTPSender(host, port, user, password)

	if err != nil {
		log.Fatal(err)
	}
	return DefaultSend(name, to, content)
}


