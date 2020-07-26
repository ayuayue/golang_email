package email

import (
	"github.com/gobuffalo/buffalo"
)
var emailclient string
func Receive(c buffalo.Context) error {
	return c.Render(200, r.HTML("email.html"))
}
func Send(c buffalo.Context) error {
	return c.Render(200, r.HTML("email.html"))

}
func SendCreateGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("email.html"))

}
func SendCreatePost(c buffalo.Context) error {
	return c.Render(200, r.HTML("email.html"))

}
