package web

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/cosban/gohst/auth"
)

type Params struct {
	Title string
	Body  []byte
}

/**
 * For backend specific pages
 */
func BackendHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]

	if !auth.IsConnected(r) {
		login := LoadPage(w, "login")
		RenderTemplate(w, login.Filename)
		return
	}

	p := LoadPage(w, title)
	switch title[len("backend/"):] {
	case "edit":
		toLoad := r.FormValue("p")
		formPage := loadPageForEdit(w, toLoad)
		p.Info = Params{formPage.Title, formPage.Body}
		log.Printf("Additional page info loaded for %s with \n%s", formPage.Info.Title, formPage.Info.Body)
		RenderTemplate(w, p.Filename)
		break
	case "manage":
		RenderTemplate(w, p.Filename)
		break
	}
}

/**
 * Loads a page from the cache, if page isn't in cache, loads the blank page
 */
func loadPageForEdit(w http.ResponseWriter, title string) *Page {
	if title == "base" || title == "blank" {
		return BlankPage(w)
	}

	filename := "./templates/" + title + ".html"
	if strings.Contains(title, "/") {
		title = title[strings.LastIndex(title, "/")+1:]
	}

	if page, ok := pages[filename]; ok {
		return page
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return BlankPage(w)
	}

	// page exists add it to the cache
	tmpl := template.Must(template.ParseFiles(filename, "templates/base.html"))
	pages[filename] = &Page{Title: title, Filename: filename, Template: tmpl, Body: body}
	return pages[filename]
}
