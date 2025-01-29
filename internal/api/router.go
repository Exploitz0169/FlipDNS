package api

import (
	"net/http"

	"github.com/exploitz0169/flipdns/internal/api/handler"
)

func loadRoutes(router *http.ServeMux) {

	handler := handler.NewHandler()

	router.HandleFunc("GET /test", handler.Test)

}
