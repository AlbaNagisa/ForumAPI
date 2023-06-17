package route

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GithubAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func GithubAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	url := "https://github.com/login/oauth/access_token?client_id=2175e85bc507a68d010a&client_secret=160b9b0ceb499b459eb6c4b417120825a94166f4&code=" + code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Accept", "application/json") // Ajoute l'en-tête "Accept" avec la valeur "application/json"

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, "Failed to retrieve access token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var githubAuthResponse GithubAuthResponse
	json.NewDecoder(resp.Body).Decode(&githubAuthResponse)
	log.Println(githubAuthResponse, resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	url = "https://api.github.com/user"
	req, err = http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", githubAuthResponse.TokenType+" "+githubAuthResponse.AccessToken)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
