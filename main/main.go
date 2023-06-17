package main

import (
	"fmt"
	"goApi/database"
	"goApi/route"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

func main() {
	r := chi.NewRouter()
	db := database.Connect()
	defer db.Close()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to my API"))
	})

	r.Post("/login", route.Login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			var id int = int(claims["id"].(float64))
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", id)))
		})
		r.Get("/me", route.Me)

		r.Post("/post", route.PostPost)
		r.Post("/upvote", route.PostVote)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", route.GetPosts)
		r.Get("/page/{page}", route.GetPostsPage)
		r.Get("/{id}", route.GetPost)
	})
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", route.GetCategories)
		r.Get("/{id}", route.GetCategory)
	})
	r.Route("/users", func(r chi.Router) {
		/* r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("All users"))
		}) */
		r.Get("/{id}", route.GetUser)

		r.Post("/", route.PostUser)
		r.Post("/discord", route.PostUserDiscord)
	})

	r.Get("/auth/github", route.GithubAuth)

	http.ListenAndServe(":8080", r)
}
