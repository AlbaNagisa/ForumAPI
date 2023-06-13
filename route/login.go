package route

import (
	"encoding/json"
	"goApi/database"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user database.User
	json.NewDecoder(r.Body).Decode(&user)
	user = database.CheckUser(user.Email, user.Password)
	var tokenString string
	if user.Email == "" {
		var tokenStruct Token
		tokenStruct.Token = "ERROR"
		jsonData, err := json.Marshal(tokenStruct)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonData)
		return
	}
	tokenString, err := generateToken(user.ID)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	var token Token
	token.Token = tokenString
	jsonData, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)

}

func generateToken(id int) (string, error) {
	// Créez une clé secrète pour la signature du token
	secretKey := []byte("secret")

	// Définissez les revendications du token (payload)
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Le token expirera après 24 heures
	}

	// Créez le token JWT en utilisant les revendications et la clé secrète
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
