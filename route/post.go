package route

import (
	"encoding/json"
	"goApi/database"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

type PostsPagination struct {
	HasNext bool            `json:"hasNext"`
	MaxPage int             `json:"maxPage"`
	Data    []database.Post `json:"data"`
}

func PostPost(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	id := claims["id"]

	newPost := database.CreatePost(r.Body, id)

	jsonData, err := json.Marshal(newPost)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetPostsPage(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	var dataToSend PostsPagination
	posts := database.GetPosts()
	pageN, _ := strconv.Atoi(page)

	posts, dataToSend.HasNext, dataToSend.MaxPage = paginate(posts, 50*pageN, 50)
	dataToSend.Data = posts

	jsonData, err := json.Marshal(dataToSend)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {

	posts := database.GetPosts()

	jsonData, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post := database.GetOnePost(id)

	jsonData, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func paginate(x []database.Post, skip int, size int) ([]database.Post, bool, int) {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}
	hasNext := true
	if end == len(x) {
		hasNext = false
	}

	maxPage := len(x) / size

	return x[skip:end], hasNext, maxPage
}
