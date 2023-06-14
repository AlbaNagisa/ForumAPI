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

func CreateUser(body io.Reader) User {
	var newUser User
	json.NewDecoder(body).Decode(&newUser)

	res, err := datab.Exec("INSERT INTO Image (path) VALUES (?)", newUser.Image)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()

	datab.Exec("INSERT INTO User (name, email, password, profileImage) VALUES (?, ?, ?, ?)", newUser.Name, newUser.Email, newUser.Password, id)

	return newUser
}

func CreatePost(body io.Reader, authorId int) Post {
	var newPost Post
	b, _ := io.ReadAll(body)
	json.Unmarshal(b, &newPost)
	fmt.Println(newPost)
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
	rows, err := datab.Query("SELECT id FROM Message")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var post Post
		rows.Scan(&post.ID)
		posts = append(posts, GetOnePost(strconv.Itoa(post.ID)))
	}
	return posts

}

func GetOnePost(id string) Post {
	var post Post
	rows, err := datab.Query("SELECT * FROM Message WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&post.ID, &post.AuthorId, &post.Date, &post.Title, &post.Content, &post.IsResponse)
	}
	post.Author = GetOneUser(strconv.Itoa(post.AuthorId))
	post.Author.Password = ""
	var prompts []string
	var images []string
	rows, err = datab.Query("SELECT prompt FROM Prompt INNER JOIN Message_Prompt ON Message_Prompt.prompt_id = Prompt.id WHERE Message_Prompt.message_id = ?", post.ID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var prompt string
		rows.Scan(&prompt)
		prompts = append(prompts, prompt)
	}
	post.Prompts = prompts

	rows, err = datab.Query("SELECT path FROM Image INNER JOIN Image_Message ON Image_Message.image_id = Image.id WHERE Image_Message.message_id = ?", post.ID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var image string
		rows.Scan(&image)
		images = append(images, image)
	}
	post.Images = images

	rows, err = datab.Query("SELECT name FROM Categories INNER JOIN Categorie_Message ON Categorie_Message.categorie_id = Categories.id WHERE Categorie_Message.message_id = ?", post.ID)
	if err != nil {
		log.Fatal(err)
	}

	var tags []string
	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	post.TagsName = tags
	return post
}

func GetOneUser(id string) User {
	var user User
	row := datab.QueryRow("SELECT * FROM User WHERE id = ?", id)
	row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ImageProfileId)
	datab.QueryRow("SELECT * FROM Image WHERE id = ?", user.ImageProfileId).Scan(&user.ImageProfileId, &user.Image)
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
	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Name)
		categories = append(categories, category)
	}
	return categories
}
