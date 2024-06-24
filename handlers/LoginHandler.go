package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"literary-lions-forum/handlers/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		login := r.FormValue("login")
		password := r.FormValue("password")
		fmt.Println("LOGIN: ", login, password)
		user, err := db.GetUserByUsernameOrEmail(dbConn, login)
		if err != nil {
			if err == sql.ErrNoRows {
				ErrorHandler(w, http.StatusBadRequest, "User not found")
			} else {
				ErrorHandler(w, http.StatusInternalServerError, "Database error")
			}
			return
		}

		
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, "Invalid password")
			return
		}

		
		setSession(user.ID, w)

		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	}
}

func setSession(userID int, w http.ResponseWriter) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    strconv.Itoa(userID),
		Expires:  expiration,
		HttpOnly: true,
		Secure:   false, 
	}
	http.SetCookie(w, &cookie)
}
func GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err 
	}
	userID, err := strconv.Atoi(cookie.Value) 
	if err != nil {
		return 0, err 
	}
	return userID, nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logout requested with method: %s", r.Method)
	if r.Method != "POST" {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false, 
	}
	http.SetCookie(w, &cookie)

	
	http.Redirect(w, r, "/", http.StatusFound)
}
