package middleware

import (
	"log/slog"
	"net/http"

	"github.com/exploitz0169/flipdns/internal/app"
)

func LoggerMiddleware(app *app.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info(
			"Got request:",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}
