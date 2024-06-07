package main

import (
	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/server"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates := template.Must(template.ParseFiles(
		filepath.Join("web/templates", "base.html"),
		filepath.Join("web/templates", tmpl),
	))
	err := templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func main() {
	database := db.InitDB()
	defer database.Close()
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	// Static file server
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))

	// Route definitions
	r.HandleFunc("/", server.MainPageHandler(database)).Methods("GET")
	r.HandleFunc("/profile", handlers.UserProfileHandler(database)).Methods("GET")
	r.HandleFunc("/create-post", handlers.PostCreateFormHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/posts", handlers.PostsHandler(database)).Methods("GET")
	r.HandleFunc("/users", handlers.UsersHandler(database)).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/add-comment", handlers.AddCommentHandler(database)).Methods("POST")
	r.HandleFunc("/like", handlers.LikeHandler(database)).Methods("POST")
	r.HandleFunc("/create-category", handlers.CreateCategoryHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/user/{id}", handlers.UserViewHandler(database)).Methods("GET")
	r.HandleFunc("/profileEdit", handlers.ProfileEditHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/updateProfile", handlers.UpdateProfileHandler(database)).Methods("POST")
	r.HandleFunc("/postView", handlers.PostViewHandler(database)).Methods("GET")
	r.HandleFunc("/editPost", handlers.EditPostHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/updatePost", handlers.UpdatePostHandler(database)).Methods("POST")
	r.HandleFunc("/like-comment", handlers.LikeCommentHandler(database)).Methods("POST")
	r.HandleFunc("/updateComment", handlers.EditCommentHandler(database)).Methods("POST")
	r.HandleFunc("/sortedPosts", handlers.CategoryPostsHandler(database)).Methods("GET")
	r.HandleFunc("/search", handlers.SearchHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/notifications", handlers.NotificationHandler(database)).Methods("GET")
	r.HandleFunc("/delete-post", handlers.EditPostHandler(database)).Methods("POST")
	r.HandleFunc("/delete-category", handlers.CreateCategoryHandler(database)).Methods("POST")
	r.HandleFunc("/mark-notifications-read", handlers.MarkNotificationsAsReadHandler(database)).Methods("POST")
	r.HandleFunc("/delete-comment", handlers.EditCommentHandler(database)).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")

	// Start server
	http.ListenAndServe("0.0.0.0:8100", r)
}
