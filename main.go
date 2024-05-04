package main

import (
	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/server"
	"net/http"
)

func main() {
	database := db.InitDB()
	defer database.Close()
	http.HandleFunc("/", server.MainPage)
	http.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePostHandler(database, w, r)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(database, w, r)
	})
	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostsHandler(database, w, r)
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.UsersHandler(database, w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(database, w, r)
	})
	http.HandleFunc("/add-comment", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddCommentHandler(database, w, r)
	})
	http.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {
		handlers.LikeHandler(database, w, r)
	})

	http.ListenAndServe("0.0.0.0:8000", nil)
}
