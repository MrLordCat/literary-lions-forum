package handlers

import (
	"literary-lions-forum/utils"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := utils.PageData{
		Title: "404 - Page Not Found",
	}
	utils.RenderTemplate(w, "404.html", data)
}
