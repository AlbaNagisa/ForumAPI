package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ImageProfileId int    `json:"image_profile_id"`
	Image          []byte `json:"image"`
}

type Post struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Date       string   `json:"date"`
	Tags       []int    `json:"tags"`
	TagsName   []string `json:"tagsName"`
	AuthorId   int      `json:"author_id"`
	Prompts    []string `json:"prompts"`
	Images     []string `json:"images"`
	IsResponse bool     `json:"is_response"`
	Author     User     `json:"author"`
	Votes      []Vote   `json:"votes"`
}

type Vote struct {
	UserId int  `json:"user_id"`
	PostId int  `json:"post_id"`
	Vote   bool `json:"vote"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var datab *sql.DB

// Connect to the database
func Connect() *sql.DB {
	wd, _ := os.Getwd()
	db, err := sql.Open("sqlite3", (wd + "/database/db.db"))
	if err != nil {
		log.Fatal(err)
	}
	datab = db
	return db
}
func CreateVote(body io.Reader) Vote {
	var vote Vote
	err := json.NewDecoder(body).Decode(&vote)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := datab.Query("SELECT * FROM Vote WHERE user_id = ? AND message_id = ?", vote.UserId, vote.PostId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	datab.Exec("DELETE FROM Vote WHERE user_id = ? AND message_id = ?", vote.UserId, vote.PostId)
	_, err = datab.Exec("INSERT INTO Vote (user_id, message_id, vote) VALUES (?, ?, ?)", vote.UserId, vote.PostId, vote.Vote)
	if err != nil {
		log.Fatal(err)

	}

	return vote

}
func CreateUser(body io.Reader) User {
	var newUser User
	json.NewDecoder(body).Decode(&newUser)
	res, err := datab.Exec("INSERT INTO Image (path) VALUES (?)", newUser.Image)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	fmt.Println(newUser.ID)
	if newUser.ID != 0 {
		datab.Exec("INSERT INTO User (name, email, password, profileImage, id) VALUES (?, ?, ?, ?, ?)", newUser.Name, newUser.Email, newUser.Password, id, newUser.ID)

	} else {
		_, err := datab.Exec("INSERT INTO User (name, email, password, profileImage) VALUES (?, ?, ?, ?)", newUser.Name, newUser.Email, newUser.Password, id)
		if err != nil {
			datab.Exec("UPDATE User SET password = ? WHERE email = ?", newUser.Password, newUser.Email)
		}

	}

	return newUser
}

func CreatePost(body io.Reader, authorId int) Post {
	var newPost Post
	b, _ := io.ReadAll(body)
	json.Unmarshal(b, &newPost)
	newPost.AuthorId = authorId

	resM, err := datab.Exec("INSERT INTO Message (title, content, date, is_response, author_id) VALUES (?, ?, ?, ?, ?)", newPost.Title, newPost.Content, newPost.Date, newPost.IsResponse, authorId)
	messId, _ := resM.LastInsertId()

	for i := 0; i < len(newPost.Images); i++ {
		resI, err := datab.Exec("INSERT INTO Image (path) VALUES (?)", newPost.Images[i])
		imageId, _ := resI.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		datab.Exec("INSERT INTO Image_Message (message_id, image_id) VALUES (?, ?)", messId, imageId)
	}

	for i := 0; i < len(newPost.Prompts); i++ {
		resP, err := datab.Exec("INSERT INTO Prompt (prompt) VALUES (?)", newPost.Prompts[i])
		if err != nil {
			log.Fatal(err)
		}
		promptId, err := resP.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		datab.Exec("INSERT INTO Message_Prompt (message_id, prompt_id) VALUES (?, ?)", messId, promptId)
	}

	for i := 0; i < len(newPost.Tags); i++ {
		datab.Exec("INSERT INTO Categorie_Message (message_id, categorie_id) VALUES (?, ?)", messId, newPost.Tags[i])
	}

	if err != nil {
		log.Fatal(err)
	}

	newPost.ID = int(messId)
	return newPost
}

func GetPosts() []Post {
	var posts []Post

	tx, err := datab.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT id FROM Message")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	postIDs := make([]string, 0) // Stocker les IDs des posts

	for rows.Next() {
		var postID string
		err := rows.Scan(&postID)
		if err != nil {
			log.Fatal(err)
		}
		postIDs = append(postIDs, postID)
	}

	// Utiliser des channels pour recevoir les résultats des goroutines
	postChan := make(chan Post)

	// Lancer les goroutines pour récupérer les posts de manière concurrente
	for _, postID := range postIDs {
		go func(id string) {
			post := GetOnePost(id)
			postChan <- post
		}(postID)
	}

	// Récupérer les résultats des channels
	for range postIDs {
		post := <-postChan
		posts = append(posts, post)
	}

	return posts
}

func GetOnePost(id string) Post {
	var post Post

	err := datab.QueryRow("SELECT * FROM Message WHERE id = ?", id).Scan(&post.ID, &post.AuthorId, &post.Date, &post.Title, &post.Content, &post.IsResponse)
	if err != nil {
		log.Fatal(err)
	}

	post.Author = GetOneUser(strconv.Itoa(post.AuthorId))
	post.Author.Password = ""

	// Utiliser des channels pour récupérer les résultats de manière asynchrone
	promptsCh := make(chan []string)
	imagesCh := make(chan []string)
	tagsCh := make(chan []string)
	voteCh := make(chan []Vote)

	// Exécuter les goroutines en parallèle pour récupérer les données
	go getPromptsForPost(id, promptsCh)
	go getImagesForPost(id, imagesCh)
	go getTagsForPost(id, tagsCh)
	go getVotesForPost(id, voteCh)

	// Récupérer les résultats des goroutines
	post.Prompts = <-promptsCh
	post.Images = <-imagesCh
	post.TagsName = <-tagsCh
	post.Votes = <-voteCh

	return post
}

func getPromptsForPost(postID string, promptsCh chan<- []string) {
	rows, err := datab.Query("SELECT prompt FROM Prompt INNER JOIN Message_Prompt ON Message_Prompt.prompt_id = Prompt.id WHERE Message_Prompt.message_id = ?", postID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var prompts []string
	for rows.Next() {
		var prompt string
		err := rows.Scan(&prompt)
		if err != nil {
			log.Fatal(err)
		}
		prompts = append(prompts, prompt)
	}

	promptsCh <- prompts
}

func getImagesForPost(postID string, imagesCh chan<- []string) {
	rows, err := datab.Query("SELECT path FROM Image INNER JOIN Image_Message ON Image_Message.image_id = Image.id WHERE Image_Message.message_id = ?", postID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var image string
		err := rows.Scan(&image)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, image)
	}

	imagesCh <- images
}

func getTagsForPost(postID string, tagsCh chan<- []string) {
	rows, err := datab.Query("SELECT name FROM Categories INNER JOIN Categorie_Message ON Categorie_Message.categorie_id = Categories.id WHERE Categorie_Message.message_id = ?", postID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, tag)
	}

	tagsCh <- tags
}

func getVotesForPost(postID string, votesCh chan<- []Vote) {
	rows, err := datab.Query("SELECT * FROM Vote WHERE message_id = ?", postID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var votes []Vote
	for rows.Next() {
		var vote Vote
		err := rows.Scan(&vote.UserId, &vote.PostId, &vote.Vote)
		if err != nil {
			log.Fatal(err)
		}
		votes = append(votes, vote)
	}

	votesCh <- votes
}

func GetOneUser(id string) User {
	var user User

	err := datab.QueryRow("SELECT User.*, Image.path FROM User INNER JOIN Image ON User.profileImage = Image.id WHERE User.id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ImageProfileId, &user.Image)
	if err != nil {
		return User{}
	}

	return user
}

func CheckUser(email string, password string) User {
	var user User
	row := datab.QueryRow("SELECT * FROM User WHERE email = ? AND password = ?", email, password)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ImageProfileId)
	if err != nil {
		return User{}
	}
	datab.QueryRow("SELECT * FROM Image WHERE id = ?", user.ImageProfileId).Scan(&user.ImageProfileId, &user.Image)
	return user
}

func GetOneCategory(id string) Category {
	var category Category
	row := datab.QueryRow("SELECT * FROM Categories WHERE id = ?", id)
	row.Scan(&category.ID, &category.Name)
	return category
}

func GetCategories() []Category {
	var categories []Category
	rows, err := datab.Query("SELECT * FROM Categories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Name)
		categories = append(categories, category)
	}
	return categories
}
