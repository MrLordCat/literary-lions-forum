package server

import (
	"database/sql"

	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"net/http"
	"text/template"
)

func MainPage(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Загрузка категорий
	categories, err := db.GetAllCategories(dbConn)
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	posts, err := db.GetAllPosts(dbConn, 0, 0)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	userID, err := handlers.GetUserIDFromSession(r)
	loggedIn := err == nil && userID != 0 // пользователь считается вошедшим в систему, если нет ошибки и userID не 0

	// Парсинг шаблона с передачей категорий и статуса аутентификации
	tmpl := template.Must(template.ParseFiles("web/template/mainPage.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"Categories": categories,
		"LoggedIn":   loggedIn, // Передаем статус аутентификации в шаблон
		"Posts":      posts,
	})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
