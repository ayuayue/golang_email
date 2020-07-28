package middlewares

import (
	// "fmt"

	"fmt"

	"github.com/gobuffalo/buffalo"
)

// LoginMiddleware 登录验证
func LoginMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler
		fmt.Println(c.Session().Get("name"))
		if c.Session().Get("name") == nil || c.Session().Get("password") == nil {
			c.Flash().Add("warning", "用户信息已过期,请重新登录")
			return c.Redirect(302, "/login")
		}
		err := next(c)
		// do some work after calling the next handler
		return err
	}
}
