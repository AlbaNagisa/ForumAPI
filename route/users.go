package route

import (
	"encoding/json"
	"goApi/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	newUser := database.CreateUser(r.Body)

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := database.GetOneUser(id)
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func PostUserDiscord(w http.ResponseWriter, r *http.Request) {
	newUser := database.CreateUser(r.Body)

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
