package middleware

import (
	"net/http"

	"github.com/exploitz0169/flipdns/internal/app"
)

type Middleware func(app *app.App, next http.Handler) http.Handler

func CreateStack(middlewares ...Middleware) Middleware {
	return func(app *app.App, next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(app, next)
		}

		return next
	}
}
