package middleware

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
)

// FireShine 自定义中间件
func FireShine() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			fmt.Printf("[FireShine] monitor, req:%+v\n", req)
			reply, err = handler(ctx, req)
			if err != nil {
				fmt.Printf("[FireShine] error, err:%+v", err)
			}
			fmt.Printf("[FireShine] reply, reply:%+v", reply)
			return

		}
	}
}
