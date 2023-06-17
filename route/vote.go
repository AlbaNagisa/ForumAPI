package route

import (
	"encoding/json"
	"fmt"
	"goApi/database"
	"net/http"
)

func PostVote(w http.ResponseWriter, r *http.Request) {
	vote := database.CreateVote(r.Body)
	jsonData, _ := json.Marshal(vote)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("vote created")
	w.Write(jsonData)
}
