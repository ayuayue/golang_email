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
	em.GET("/send", email.Receive)
	em.GET("/send/index", email.Send)
	em.GET("/send/create", email.SendCreateGet)
	em.POST("/send/create", email.SendCreatePost)
}
