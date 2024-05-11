package handlers

import (
	"database/sql"
	"literary-lions-forum/handlers/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	user, err := db.GetUserByUsernameOrEmail(dbConn, login)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusBadRequest)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Сравнение предоставленного пароля с хэшированным паролем в базе данных
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Установка сессии для пользователя после успешного входа
	setSession(user.ID, w)

	// Перенаправление пользователя на главную страницу или страницу профиля
	http.Redirect(w, r, "/", http.StatusFound)
}

func setSession(userID int, w http.ResponseWriter) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    strconv.Itoa(userID),
		Expires:  expiration,
		HttpOnly: true,
		Secure:   false, // Этот флаг следует включать, если вы используете HTTPS
	}
	http.SetCookie(w, &cookie)
}
func GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err // Ошибка, если куки нет
	}
	userID, err := strconv.Atoi(cookie.Value) // Преобразование значения куки обратно в int
	if err != nil {
		return 0, err // Ошибка, если значение не является целым числом
	}
	return userID, nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logout requested with method: %s", r.Method)
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Удаление сессионной куки
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false, // Включите для HTTPS
	}
	http.SetCookie(w, &cookie)

	// Перенаправление на страницу входа или на главную страницу
	http.Redirect(w, r, "/", http.StatusFound)
}
