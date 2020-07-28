package routes

import (
	"eclient/actions/email"

	"github.com/gobuffalo/buffalo"
)

// WebRoutes web路由
func WebRoutes(app *buffalo.App) {
	em := app.Group("/mails")
	em.GET("/", email.Receive)
	em.GET("/index", email.Receive)
	em.GET("/send", email.Send)
	em.POST("/send", email.SendAct)
}