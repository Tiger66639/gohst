package web

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/cosban/gohst/auth"
	"github.com/cosban/gohst/data"
)

const backendLocation = "templates/"

// EditParams is a struct which holds the info needed to load a page into the
// page editor
type EditParams struct {
	Title, Body string
}

type ManageParams struct {
	Pages []*data.Page
}

// BackendHandler is used to route backend specific pages
func BackendHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	log.Printf(title)

	if !auth.IsConnected(r) {
		login := LoadPage(w, PageLocation, "login")
		RenderTemplate(w, login)
		return
	}

	switch title[len("backend/"):] {
	case "edit":
		edit(w, r, title)
	case "edit/submit":
		submitEdit(w, r, title)
	case "":
		title = "backend/manage"
		fallthrough
	case "manage":
		manage(w, r, title)
	default:
		p := LoadPage(w, backendLocation, title)
		p.LoggedIn = true
		RenderTemplate(w, p)
	}
}

func edit(w http.ResponseWriter, r *http.Request, title string) {
	p := LoadPage(w, backendLocation, title)
	toLoad := r.FormValue("p")
	formPage := loadPageForEdit(w, toLoad)
	p.Info = EditParams{toLoad, string(formPage.Body)}
	p.LoggedIn = true
	RenderTemplate(w, p)
}

func submitEdit(w http.ResponseWriter, r *http.Request, title string) {
	formURL := r.FormValue("p")
	formPage := loadPageForEdit(w, formURL)
	formPage.Body = []byte(r.FormValue("body"))
	formPage.Template = template.Must(template.New(title).Parse(string(formPage.Body)))
	formPage.Template = template.Must(formPage.Template.ParseFiles("templates/shared/base.html"))
	pages[title] = formPage
	http.Redirect(w, r, "/backend/manage", http.StatusFound)
}

func manage(w http.ResponseWriter, r *http.Request, title string) {
	p := LoadPage(w, backendLocation, title)
	pages, err := data.GetAllPages()
	if err != nil {
		panic(err)
	}
	p.Info = ManageParams{pages}
	p.LoggedIn = true
	RenderTemplate(w, p)
}

/**
 * Loads a page from the cache, if page isn't in cache, loads the blank page
 */
func loadPageForEdit(w http.ResponseWriter, title string) *Page {
	if title == "base" || title == "blank" {
		return BlankPage(w)
	}

	filename := PageLocation + title + ".html"
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
	tmpl := template.Must(template.ParseFiles(filename, SharedLocation+"/base.html"))
	pages[filename] = &Page{Title: title, Template: tmpl, Body: body}
	return pages[filename]
}
