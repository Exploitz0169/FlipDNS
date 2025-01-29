package api

import (
	"log/slog"
	"net/http"

	"github.com/exploitz0169/flipdns/internal/api/middleware"
	"github.com/exploitz0169/flipdns/internal/app"
)

type API struct {
	app *app.App
}

func NewAPI(app *app.App) *API {
	return &API{
		app: app,
	}
}

func (a *API) Run() {

	router := http.NewServeMux()
	loadRoutes(router)

	middlewareStack := middleware.CreateStack(
		middleware.LoggerMiddleware,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: middlewareStack(a.app, router),
	}

	a.app.Logger.Info("API server starting", slog.String("addr", server.Addr))
	err := server.ListenAndServe()
	if err != nil {
		a.app.Logger.Error(err.Error())
	}

}
