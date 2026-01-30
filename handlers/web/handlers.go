package web

import (
	"net/http"
	"strings"
	"web-demo/server"
)

type WebHandler struct {
	App *server.Application
}

func NewWebHandler(app *server.Application) *WebHandler {
	return &WebHandler{App: app}
}

// --- Page Handlers ---

func (h *WebHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path != "/" {
		h.NotFoundHandler(w, r)
		return
	}

	h.Render(w, http.StatusOK, "index.html")
}

func (h *WebHandler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	h.Render(w, http.StatusOK, "about.html")
}

func (h *WebHandler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	h.Render(w, http.StatusNotFound, "404.html")
}
