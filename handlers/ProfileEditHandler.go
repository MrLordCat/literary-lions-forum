package handlers

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"net/http"
)

func ProfileEditHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "You need to be logged in to access this page", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(dbConn, userID)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/template/profileEdit.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"User": user,
	})
	if err != nil {
		http.Error(w, "Failed to render profile edit page", http.StatusInternalServerError)
	}
}

func UpdateProfileHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	if !db.СheckCurrentPassword(dbConn, userID, currentPassword) {
		http.Error(w, "Current password is incorrect", http.StatusForbidden)
		return
	}

	// Проверка, что новый пароль и подтверждение совпадают
	if newPassword != confirmPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Обновление профиля и пароля, если все проверки прошли успешно
	if err = db.UpdateUser(dbConn, userID, username, firstName, lastName, newPassword); err != nil {
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusFound)
}
