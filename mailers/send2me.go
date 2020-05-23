package mailers

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
)

func SendSend2mes(from string, desc string) error {
	m := mail.NewMessage()

	// fill in with your stuff:
	m.Subject = "从博客发送过来"
	m.From = "2294595856@qq.com"
	m.To = []string{"1401262639@qq.com"}
	content := "来自: " + from + " ,内容为: " + desc
	err := m.AddBody(r.String(content), render.Data{})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
