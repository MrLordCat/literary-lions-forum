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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.MainPage(database, w, r)
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.UserProfileHandler(database, w, r)
	})
	http.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostCreateFormHandler(database, w, r)
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
	http.HandleFunc("/create-category", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCategoryHandler(database, w, r)
	})
	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		handlers.UserViewHandler(database, w, r)
	})
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.ListenAndServe("0.0.0.0:8100", nil)
}
