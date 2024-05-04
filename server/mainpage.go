package server

import (
	"net/http"
	"text/template"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/mainPage.html"))
	tmpl.Execute(w, tmpl)
}
