package server

import (
	"encoding/json"
	"net/http"
)

// --- API Handlers ---

func (app *Application) EchoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"data": "hello world"}
	json.NewEncoder(w).Encode(data)
}
