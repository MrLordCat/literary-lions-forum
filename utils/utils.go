// utils/templates.go
package utils

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	funcMap := template.FuncMap{
		"renderPostContent": RenderPostContent,
	}

	templates := template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
		filepath.Join("web/templates/", "base.html"),
		filepath.Join("web/templates/", tmpl),
	))

	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PageData struct {
	Title               string
	User                db.User
	LoggedIn            bool
	IsAdmin             bool
	Karma               int
	UserPosts           []db.Post
	LikedPosts          []db.Post
	Notifications       []db.Notification
	UnreadNotifications int
	Posts               []db.Post
	Categories          []db.Category
	Users               []db.User
	IsOwnProfile        bool
}

func GetPageData(dbConn *sql.DB, userID int, options map[string]bool) (PageData, error) {
	var data PageData
	var err error

	if userID > 0 {
		data.User, err = db.GetUserByID(dbConn, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Если пользователь не найден, продолжим выполнение без установки LoggedIn
				data.LoggedIn = false
			} else {
				return data, err
			}
		} else {
			data.LoggedIn = true
			data.IsAdmin = data.User.IsAdmin
		}
		if options["karma"] {
			data.Karma = data.User.Karma
		}

		if options["userPosts"] {
			data.UserPosts, err = db.GetUserPosts(dbConn, userID)
			if err != nil {
				return data, err
			}
		}

		if options["likedPosts"] {
			data.LikedPosts, err = db.GetLikedPosts(dbConn, userID)
			if err != nil {
				return data, err
			}
		}

		if options["notifications"] {
			data.Notifications, err = db.GetUnreadNotifications(dbConn, userID)
			if err != nil {
				return data, err
			}
			data.UnreadNotifications = len(data.Notifications)
		}
	}

	if options["posts"] {
		data.Posts, err = db.GetAllPosts(dbConn, 0, 0)
		if err != nil {
			return data, err
		}
	}

	if options["categories"] {
		data.Categories, err = db.GetAllCategories(dbConn)
		if err != nil {
			return data, err
		}
	}
	if options["isOwnProfile"] {
		data.IsOwnProfile = options["isOwnProfile"]
	}

	return data, nil
}
func RenderPostContent(content string) template.HTML {
	extensions := blackfriday.NoIntraEmphasis |
		blackfriday.Tables |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.BackslashLineBreak |
		blackfriday.DefinitionLists |
		blackfriday.AutoHeadingIDs
	content = strings.ReplaceAll(content, "\n", "<br>")
	output := blackfriday.Run([]byte(content), blackfriday.WithExtensions(extensions))
	return template.HTML(output)
}
