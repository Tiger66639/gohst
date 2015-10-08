package web

import (
	"io/ioutil"
	"net/http"
	"text/template"
)

func SubmitEdit(w http.ResponseWriter, r *http.Request) {
	formURL := r.FormValue("p")
	formPage := loadPageForEdit(w, formURL)
	formPage.Body = []byte(r.FormValue("body"))
	formPage.Template = template.Must(template.New(formPage.Title).Parse(string(formPage.Body)))
	formPage.Template = template.Must(formPage.Template.ParseFiles("templates/shared/base.html"))
	pages[formPage.Title] = formPage
	ioutil.WriteFile("./"+PageLocation+formPage.Title+".html", formPage.Body, 0644)
	http.Redirect(w, r, "/backend/manage", http.StatusFound)
}
