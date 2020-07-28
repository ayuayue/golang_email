package email

import (
	"eclient/mailers"
	"errors"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var host = envy.Get("IMAP_HOST", "")
var port = envy.Get("IMAP_PORT", "")
var server = fmt.Sprintf("%s:%s", host, port)

type Msg struct {
	ID      uint32 `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

func Receive(c buffalo.Context) error {
	var mailItems []*mailers.MailItem
	msg := Msg{}
	msgs := []Msg{}
	name := c.Session().Get("name")
	password := c.Session().Get("password")
	mailItems = mailers.GerReceiveMail(host, port, name.(string), password.(string))
	if mailItems == nil || len(mailItems) == 0 {
		return errors.New("未接收到邮件，请检查邮箱配置或邮箱为空")
	}
	// 获取邮件内容，并将接收到的邮件存进库中
	for _, receive := range mailItems {
		message := mailers.GetReceiveMailMessage(host, port, name.(string), password.(string), receive.ID)
		msg.ID = receive.ID
		msg.From = receive.From
		msg.To = name.(string)
		msg.Subject = message.Subject
		msg.Content = message.Body
		msg.Date = message.Date
		msgs = append(msgs, msg)
	}
	c.Set("mails", msgs)
	return c.Render(200,r.HTML("/mails/index"))

}
func Send(c buffalo.Context) error {
	return c.Render(200, r.HTML("mails/index.html"))

}
func SendAct(c buffalo.Context) error {
	fmt.Println(c.Params())
	c.Flash().Add("success", "邮件发送成功")
	return c.Redirect(302, "/mails/")
}
