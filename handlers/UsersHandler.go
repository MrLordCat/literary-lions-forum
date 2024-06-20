package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
	"sort"
)

func UsersHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID пользователя из сессии, если вход выполнен
		userID, _ := GetUserIDFromSession(r)

		users, err := db.GetAllUsers(dbConn)
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		// Сортируем пользователей по карме по умолчанию
		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" || sortBy == "karma" {
			sort.Slice(users, func(i, j int) bool {
				karmaI := int64(0)
				if users[i].Karma.Valid {
					karmaI = users[i].Karma.Int64
				}
				karmaJ := int64(0)
				if users[j].Karma.Valid {
					karmaJ = users[j].Karma.Int64
				}
				return karmaI > karmaJ
			})
		}

		// Получаем топ 10 пользователей по карме
		topUsers := users
		if len(users) > 10 {
			topUsers = users[:10]
		}

		options := map[string]bool{
			"notifications": userID != 0,
		}
		pageData, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pageData.Title = "Users"
		pageData.Users = users
		pageData.TopUsers = topUsers // Добавляем топ пользователей

		err = utils.RenderTemplate(w, "usersList.html", pageData)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}
