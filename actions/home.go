package actions

import (
	"eclient/mailers"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	// return c.Redirect(302, "emailPath()")
	return c.Render(200, r.HTML("index.html"))
}

func Login(c buffalo.Context) error {
	c.Set("name", "")
	c.Set("password", "")
	return c.Render(200, r.HTML("login.html"))
}

func LoginAction(c buffalo.Context) error {
	name := c.Param("name")
	password := c.Param("password")
	fmt.Println(name, password)
	return c.Redirect(302, "/login")
}

func Nologin(c buffalo.Context) error {
	c.Set("name", "")
	c.Set("to", "1401262639@qq.com")
	c.Set("content", "")
	return c.Render(200, r.HTML("nologin.html"))
}

func NologinAction(c buffalo.Context) error {
	name := c.Param("name")
	to := c.Param("to")
	content := c.Param("content")
	vrs := validate.Validate(
		&validators.EmailIsPresent{Field: name, Name: "您的邮箱", Message: " .不合法. "},
		&validators.EmailLike{Field: name, Name: "您的邮箱", Message: " .格式不正确. "},
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
