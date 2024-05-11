package server

import (
	"database/sql"

	"literary-lions-forum/handlers"
	"literary-lions-forum/handlers/db"
	"net/http"
	"text/template"
)

func MainPage(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	userID, err := handlers.GetUserIDFromSession(r)
	loggedIn := err == nil && userID != 0
	//fmt.Println("Logged in user ID:", userID)

	// Получаем посты
	posts, err := db.GetAllPosts(dbConn, 0, 0)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return // После установки HTTP ошибки функция должна завершиться
	}

	// Получаем категории
	categories, err := db.GetAllCategories(dbConn)
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return // Также выходим, если не удалось получить категории
	}

	// Парсим шаблон и передаем данные
	tmpl, err := template.ParseFiles("web/template/mainPage.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return // Обрабатываем возможную ошибку парсинга шаблона
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"Categories": categories,
		"LoggedIn":   loggedIn,
		"Posts":      posts,
		"UserID":     userID,
	})

	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return // После ошибки рендеринга, функция также должна завершиться
	}
}
