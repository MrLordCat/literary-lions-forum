package main

import (
	"encoding/json"
	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/server"
	"log"
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
	r.Use(LoggingMiddleware)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("web/static"))))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir("uploads"))))

	r.HandleFunc("/", server.MainPageHandler(database)).Methods("GET")
	r.HandleFunc("/profile", handlers.UserProfileHandler(database)).Methods("GET")
	r.HandleFunc("/create-post", handlers.PostCreateFormHandler(database)).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler(database)).Methods("GET", "POST")

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
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	http.ListenAndServe("0.0.0.0:8000", r)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(lw, r)

		log.Printf("[%s] %s %d", r.Method, r.URL.Path, lw.statusCode)
		if lw.data != nil {
			logData, _ := json.MarshalIndent(lw.data, "", "  ")
			log.Printf("Data: %s", logData)
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	data       interface{}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {

	if json.Unmarshal(data, &lrw.data) != nil {
		lrw.data = nil
	}
	return lrw.ResponseWriter.Write(data)
}
