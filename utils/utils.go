// utils/templates.go
package utils

import (
	"database/sql"
	"html/template"
	"literary-lions-forum/handlers/db"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	funcMap := template.FuncMap{
		"renderPostContent": RenderPostContent,
	}

	templates, err := template.New("base.html").Funcs(funcMap).ParseFiles(
		filepath.Join("web/templates/", "base.html"),
		filepath.Join("web/templates/home/", "catAdmin.html"),
		filepath.Join("web/templates/home/", "categoriesBlock.html"),
		filepath.Join("web/templates/home/", "postsBlock.html"),
		filepath.Join("web/templates/home/", "top-usersBlock.html"),
		filepath.Join("web/templates/", "notifications.html"),
		filepath.Join("web/templates/", "usersList.html"),
		filepath.Join("web/templates/", tmpl),
	)
	if err != nil {
		return err
	}

	err = templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		return err
	}

	return nil
}

type PageData struct {
	Title               string
	User                db.User
	LoggedIn            bool
	IsAdmin             bool
	Karma               sql.NullInt64
	UserPosts           []db.Post
	LikedPosts          []db.Post
	Notifications       []db.Notification
	UnreadNotifications int
	Posts               []db.Post
	Categories          []db.Category
	Users               []db.User
	TopUsers            []db.User
	IsOwnProfile        bool
	IsProfile           bool
	SinglePost          db.Post
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
			data.UserPosts, err = db.GetAllPosts(dbConn, 0, int64(userID))
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

	if options["topUsers"] {
		users, err := db.GetAllUsers(dbConn)
		if err != nil {
			return data, err
		}
		// Сортируем пользователей по карме
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

		// Получаем топ 10 пользователей
		if len(users) > 10 {
			data.TopUsers = users[:10]
		} else {
			data.TopUsers = users
		}
	}

	if options["isOwnProfile"] {
		data.IsOwnProfile = options["isOwnProfile"]
	}
	var postID int
	if options["singlePost"] {
		posts, err := db.GetAllPosts(dbConn, postID, 0)
		if err != nil {
			return data, err
		}
		if len(posts) > 0 {
			data.SinglePost = posts[0]
		}
	}
	// Установка значения IsProfile
	if options["isProfile"] {
		data.IsProfile = true
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
