package actions

import (
	"eclient/mailers"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	// return c.Render(200,r.JSON(c))
	return c.Render(200, r.HTML("index.html"))
}
// Login 登录页
func Login(c buffalo.Context) error {
	c.Set("name", "")
	c.Set("password", "")
	return c.Render(200, r.HTML("login.html"))
}
// LoginAction 登录逻辑处理
func LoginAction(c buffalo.Context) error {
	name := c.Param("name")
	password := c.Param("password")
	host := envy.Get("IMAP_HOST", "localhost")
	port := envy.Get("IMAP_PORT", "")
	fmt.Println(name, password)
	vrs := validate.Validate(
		&validators.EmailIsPresent{Field: name, Name: "您的邮箱", Message: " .不合法. "},
		&validators.EmailLike{Field: name, Name: "您的邮箱", Message: " .格式不正确. "},
		&validators.StringIsPresent{Field: password, Name: "密码", Message: "不能为空"},
	)
	if vrs.HasAny() {
		c.Set("errors", vrs.Errors)
		c.Set("name", name)
		c.Set("password", password)
		return c.Render(200, r.HTML("login.html"))
	}
	ok := mailers.Connect(fmt.Sprintf("%s:%s", host, port), name, password)
	if !ok {
		c.Set("name", name)
		c.Set("password", password)
		c.Flash().Add("danger", fmt.Sprintf("%s", "登录失败"))
		return c.Render(200, r.HTML("login.html"))
	}
	c.Session().Set("name", name)
	c.Session().Set("password", password)
	c.Flash().Add("success", "登录成功")
	return c.Redirect(302, "/mails")
}
// Nologin 免登录页
func Nologin(c buffalo.Context) error {
	c.Set("name", "")
	c.Set("to", "1401262639@qq.com")
	c.Set("content", "")
	return c.Render(200, r.HTML("nologin.html"))
}
// NologinAction 免登录逻辑处理
func NologinAction(c buffalo.Context) error {
	name := c.Param("name")
	to := c.Param("to")
	content := c.Param("content")
	vrs := validate.Validate(
		&validators.EmailIsPresent{Field: name, Name: "您的邮箱", Message: " .不合法. "},
		&validators.EmailLike{Field: name, Name: "您的邮箱", Message: " .格式不正确. "},
		&validators.StringIsPresent{Field: content, Name: "内容", Message: "不能为空"},

		&validators.EmailIsPresent{Field: to, Name: "收件人邮箱", Message: " .不合法. "},
		&validators.EmailLike{Field: to, Name: "收件人邮箱", Message: " .格式不正确. "},
	)
	if vrs.HasAny() {
		c.Set("errors", vrs.Errors)
		c.Set("name", name)
		c.Set("to", to)
		c.Set("content", content)
		return c.Render(200, r.HTML("nologin.html"))
	}
	err := mailers.DefaultSend(name, to, content)
	if err != nil {
		c.Set("name", name)
		c.Set("to", to)
		c.Set("content", content)
		c.Flash().Add("danger", fmt.Sprintf("%s", err))
		return c.Render(200, r.HTML("nologin.html"))

	}
	c.Flash().Add("success", "邮件已发送!")
	return c.Redirect(302, "/login")
}
// Logout 退出
func Logout(c buffalo.Context) error {
	c.Flash().Clear()
	c.Session().Clear()
	return c.Redirect(302, "/login")
}