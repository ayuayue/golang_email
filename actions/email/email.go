package email

import (
	"eclient/mailers"
	"errors"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

var host = envy.Get("IMAP_HOST", "")
var port = envy.Get("IMAP_PORT", "")
var server = fmt.Sprintf("%s:%s", host, port)

// Msg 邮件消息结构体
type Msg struct {
	ID      uint32 `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

var mailItems []*mailers.MailItem
var msgs []Msg

// Receive 接受邮件处理
func Receive(c buffalo.Context) error {
	msg := Msg{}
	name := c.Session().Get("name")
	password := c.Session().Get("password")
	msgAmount := c.Session().Get("msg_amount")
	mailAmount := mailers.GetMailNum(server, name.(string), password.(string))["INBOX"]
	if msgAmount != nil && msgAmount.(int) < mailAmount {
		mailItems = mailers.GerReceiveMail(host, port, name.(string), password.(string))
		msgs = []Msg{}
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
	}
	c.Set("mails", msgs)
	c.Session().Set("msg_amount", len(msgs))
	return c.Render(200, r.HTML("/mails/index"))

}

// Send 发送邮件页
func Send(c buffalo.Context) error {
	c.Set("content", "")
	c.Set("to", "")
	return c.Render(200, r.HTML("/mails/new.plush.html"))

}

// SendAct 处理邮件发送
func SendAct(c buffalo.Context) error {
	to := c.Param("to")
	content := c.Param("content")
	name := c.Session().Get("name")
	password := c.Session().Get("password")
	vrs := validate.Validate(
		&validators.EmailIsPresent{Field: to, Name: "邮箱", Message: " .不合法. "},
		&validators.EmailLike{Field: to, Name: "邮箱", Message: " .格式不正确. "},
		&validators.StringIsPresent{Field: content, Name: "内容", Message: "不能为空"},
	)
	if vrs.HasAny() {
		c.Set("errors", vrs.Errors)
		c.Set("to", to)
		c.Set("content", content)
		return c.Render(200, r.HTML("/mails/new.plush.html"))
	}
	err := mailers.SendByLogin(name.(string), password.(string), to, content)
	if err != nil {
		c.Flash().Add("danger", "邮件发送失败")
	}
	c.Flash().Add("success", "邮件发送成功")
	return c.Redirect(302, "/mails/")
}
