package web

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/cosban/gohst/auth"
)

// EditParams is a struct which holds the info needed to load a page into the
// page editor
type EditParams struct {
	Title string
	Body  []byte
}

// BackendHandler is used to route backend specific pages
func BackendHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	log.Printf(title)

	if !auth.IsConnected(r) {
		login := LoadPage(w, "backend", "login")
		RenderTemplate(w, login)
		return
	}

	switch title[len("backend/"):] {
	case "edit":
		p := LoadPage(w, "backend", title)
		toLoad := r.FormValue("p")
		formPage := loadPageForEdit(w, toLoad)
		p.Info = EditParams{toLoad, formPage.Body}
		RenderTemplate(w, p)
		break
	case "edit/submit":
		formURL := r.FormValue("p")
		formPage := loadPageForEdit(w, formURL)
		formPage.Body = []byte(r.FormValue("body"))
		formPage.Template = template.Must(template.New(title).Parse(string(formPage.Body)))
		formPage.Template = template.Must(formPage.Template.ParseFiles("templates/base.html"))
		pages[title] = formPage
		http.Redirect(w, r, "/backend/manage", http.StatusFound)
		log.Printf("SECOND\n%s", formPage.Template)
		break
	case "manage":
		p := LoadPage(w, "backend", title)
		RenderTemplate(w, p)
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
	pages[filename] = &Page{Title: title, Template: tmpl, Body: body}
	return pages[filename]
}
