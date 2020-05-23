package actions

import (
	"email/mailers"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"strconv"
	"time"
)

func EmailHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("email.html"))
}

func EmailPostHandler(c buffalo.Context) error {
	fmt.Println("start handler")
	from := c.Param("email")
	desc := c.Param("content")
	per := c.Param("next_time")
	perInt, _ := strconv.ParseInt(per, 10, 64)

	now := time.Now().Unix()
	diff := now - perInt
	if diff-90 < 0 {
		c.Set("next_time", time.Now().Unix())
		c.Flash().Add("danger", "发送频繁导致失败,请间隔一分半再发送")
		return c.Render(201, r.HTML("email.html"))
	}
	c.Set("next_time", time.Now().Unix())

	mailers.SendSend2mes(from, desc)
	c.Flash().Add("success", "发送成功")
	return c.Render(200, r.HTML("email.html"))
}
