
package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"literary-lions-forum/handlers/db"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"github.com/russross/blackfriday/v2"
)

func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call: missing key or value")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
func timeSince(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := hours / 24

	if hours < 24 {
		if hours > 0 {
			return fmt.Sprintf("%d hours ago", hours)
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	return fmt.Sprintf("%d days ago", days)
}
func RenderTemplate(w http.ResponseWriter, tmpl string, data PageData) error {
	log.Println("Starting RenderTemplate")

	funcMap := template.FuncMap{
		"timeSince":         timeSince,
		"renderPostContent": RenderPostContent,
		"dict":              Dict,
	}

	templates, err := template.New("base.html").Funcs(funcMap).ParseFiles(
		filepath.Join("web/templates/", "base.html"),
		filepath.Join("web/templates/home/", "catAdmin.html"),
		filepath.Join("web/templates/home/", "categoriesBlock.html"),
		filepath.Join("web/templates/post/", "posts.html"),
		filepath.Join("web/templates/post/", "post.html"),
		filepath.Join("web/templates/post/", "postActions.html"),
		filepath.Join("web/templates/post/", "comments.html"),
		filepath.Join("web/templates/home/", "top-usersBlock.html"),
		filepath.Join("web/templates/", "notifications.html"),
		filepath.Join("web/templates/", "usersList.html"),
		filepath.Join("web/templates/", "404.html"),
		filepath.Join("web/templates/", tmpl),
	)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		return err
	}

	log.Println("Templates parsed successfully")

	err = templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}

	log.Println("Template executed successfully")
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
	Post                db.Post
	CategoryName        string
	CanEdit             bool
	IsDeleted           bool
	Sort                string
	CategoryID          int
	SearchResults       struct { 
		Posts []db.Post
		Users []db.User
	}
}

func GetPageData(dbConn *sql.DB, userID int, options map[string]bool) (PageData, error) {
	var data PageData
	var err error

	if userID > 0 {
		data.User, err = db.GetUserByID(dbConn, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				
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
			data.UserPosts, err = db.GetAllPosts(dbConn, 0, int64(userID), "likes")
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
		data.Posts, err = db.GetAllPosts(dbConn, 0, 0, "p.created_at DESC")
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
		posts, err := db.GetAllPosts(dbConn, postID, 0, "likes")
		if err != nil {
			return data, err
		}
		if len(posts) > 0 {
			data.Post = posts[0]
		}
	}
	
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

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.UseXHTML |
			blackfriday.Smartypants |
			blackfriday.SmartypantsFractions |
			blackfriday.SmartypantsDashes |
			blackfriday.SmartypantsLatexDashes,
	})

	output := blackfriday.Run([]byte(content), blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extensions))

	
	htmlContent := fmt.Sprintf("<div style=\"white-space: pre-wrap;\">%s</div>", output)
	return template.HTML(htmlContent)
}
