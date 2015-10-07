package web

import (
	"io/ioutil"
	"net/http"
	"text/template"
)

func submitEdit(w http.ResponseWriter, r *http.Request, title string) {
	formURL := r.FormValue("p")
	formPage := loadPageForEdit(w, formURL)
	formPage.Body = []byte(r.FormValue("body"))
	formPage.Template = template.Must(template.New(title).Parse(string(formPage.Body)))
	formPage.Template = template.Must(formPage.Template.ParseFiles("templates/shared/base.html"))
	pages[formPage.Title] = formPage
	ioutil.WriteFile("./"+PageLocation+formPage.Title+".html", formPage.Body, 0644)
	http.Redirect(w, r, "/backend/manage", http.StatusFound)
}
