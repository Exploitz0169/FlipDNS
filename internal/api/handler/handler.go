package handler

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) LoadRoutes(router *http.ServeMux) {
	router.HandleFunc("/test", h.Test)
}
