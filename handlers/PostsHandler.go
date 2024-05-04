package handlers

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"net/http"
)

func PostsHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetAllPosts(dbConn)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/template/posts.html"))
	err = tmpl.Execute(w, map[string]interface{}{"Posts": posts})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
