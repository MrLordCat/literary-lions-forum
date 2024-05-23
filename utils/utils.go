// utils/templates.go
package utils

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates := template.Must(template.ParseFiles(
		filepath.Join("web/templates", "base.html"),
		filepath.Join("web/templates", tmpl),
	))
	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
