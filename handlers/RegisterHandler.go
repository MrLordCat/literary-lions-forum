package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"literary-lions-forum/handlers/db"
	"net/http"
	"net/mail"
	"strings"
)

func RegisterHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println("RECEIVED DATA:", username, email, password)
		// Проверка формата email

		_, err := mail.ParseAddress(email)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, "Invalid email format")
			return
		}

		// Устанавливаем isAdmin в true, если username равен 'admin'
		isAdmin := username == "admin"

		if err := db.CreateUser(dbConn, username, email, password, isAdmin); err != nil {
			// Если пользователь уже существует, вернем соответствующую ошибку
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				ErrorHandler(w, http.StatusConflict, "Username or email already exists")
				return
			}
			ErrorHandler(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
			return
		}

		// Если регистрация успешна, отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"success": "User created successfully"})
	}
}
