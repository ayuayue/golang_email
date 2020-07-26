package middlewares

import (
	// "fmt"
	"net/url"
	"strings"

	"github.com/gobuffalo/buffalo"
)

// InputFilterMiddleware 输入过滤器
func InputFilterMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler
		//获取所有参数
		//遍历参数,对每个参数进行头尾多余空格的去除处理
		//重新放回到param
		if p, ok := c.Params().(url.Values); ok {
			for k, v := range p {
				p.Set(k, strings.Trim(v[0], " \t\n\r\x0B"))
			}
		}
		err := next(c)
		// do some work after calling the next handler
		return err
	}
}
