package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"literary-lions-forum/utils"
	"net/http"
)

func ProfileEditHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "You need to be logged in to access this page", http.StatusUnauthorized)
			return
		}

		options := map[string]bool{
			"notifications": true,
		}

		data, err := utils.GetPageData(dbConn, userID, options)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data.Title = "Edit Profile"

		utils.RenderTemplate(w, "profile/profileEdit.html", data)
	}
}
func UpdateProfileHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "You need to be logged in to update your profile", http.StatusUnauthorized)
			return
		}

		username := r.FormValue("username")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		currentPassword := r.FormValue("currentPassword")
		newPassword := r.FormValue("newPassword")
		confirmPassword := r.FormValue("confirmPassword")

		// Проверка текущего пароля
		if !db.CheckCurrentPassword(dbConn, userID, currentPassword) {
			http.Error(w, "Current password is incorrect", http.StatusForbidden)
			return
		}

		// Проверка, что новый пароль и подтверждение совпадают
		if newPassword != "" && newPassword != confirmPassword {
			http.Error(w, "New passwords do not match", http.StatusBadRequest)
			return
		}

		// Обновление профиля пользователя
		if err := db.UpdateUser(dbConn, userID, username, firstName, lastName, newPassword); err != nil {
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}
