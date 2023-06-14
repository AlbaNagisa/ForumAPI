package route

import (
	"encoding/json"
	"goApi/database"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth"
)

func Me(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var id int = int(claims["id"].(float64))
	user := database.GetOneUser(strconv.Itoa(id))
	user.Password = ""
	jsonData, err := json.Marshal(user)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Write(jsonData)
}
