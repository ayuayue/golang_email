package routes

import (
	"eclient/actions/email"
	"eclient/actions/middlewares"

	"github.com/gobuffalo/buffalo"
)

// WebRoutes web路由
func WebRoutes(app *buffalo.App) {
	app.GET("/",email.Receive)
	em := app.Group("/mails")
	em.Use(middlewares.LoginMiddleware)

	em.GET("/", email.Receive)
	em.GET("/index", email.Receive)
	em.GET("/send", email.Send)
	em.POST("/send", email.SendAct)
}
