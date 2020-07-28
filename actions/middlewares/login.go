package middlewares

import (
	// "fmt"

	"github.com/gobuffalo/buffalo"
)

// LoginMiddleware 登录验证
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler

		if c.Session().Get("name") == "" || c.Session().Get("password") == "" {
			c.Flash().Add("warning", "用户信息已过期,请重新登录")
		}
		err := next(c)
		// do some work after calling the next handler
		return err
	}
}
