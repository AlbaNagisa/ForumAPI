package route

import (
	"encoding/json"
	"goApi/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := database.GetCategories()
	jsonData, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category := database.GetOneCategory(id)
	jsonData, err := json.Marshal(category)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
