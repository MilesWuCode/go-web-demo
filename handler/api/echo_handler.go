package api

import "net/http"

// EchoHandler is a simple handler that echoes a message.
func (h *APIHandler) EchoHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"data": "hello world"}
	h.App.WriteJSON(w, http.StatusOK, data)
}
