package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse the layout and the specific template
	layoutPath := filepath.Join("web", "templates", "layouts", "base.html")
	pagePath := filepath.Join("web", "templates", "pages", tmpl+".html")

	t, err := template.ParseFiles(layoutPath, pagePath)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
