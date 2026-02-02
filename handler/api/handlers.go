package api

import "web-demo/server"

type APIHandler struct {
	App *server.Application
}

func NewAPIHandler(app *server.Application) *APIHandler {
	return &APIHandler{App: app}
}
