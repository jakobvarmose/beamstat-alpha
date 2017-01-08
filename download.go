package main

import (
	"html/template"
	"net/http"
)

func handleDownload(w http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/download.html"))
	tmpl.Execute(w, nil)
}
