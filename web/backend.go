package web

import (
	"io/ioutil"
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

type UserlistParams struct {
	Users []*data.User
}

var routes = map[string]func(http.ResponseWriter, *http.Request, string){
	"users":       users,
	"edit":        edit,
	"edit/submit": submitEdit,
	"manage":      manage,
	"":            users,
}

// BackendHandler is used to route backend specific pages
func BackendHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]

	if !auth.IsConnected(r) {
		login := LoadPage(w, PageLocation, "login")
		RenderTemplate(w, r, login)
		return
	}

	if route, ok := routes[title[len("backend/"):]]; ok {
		if len(title[len("backend/"):]) == 0 {
			title = "backend/users"
		}
		route(w, r, title)
	} else {
		p := LoadPage(w, backendLocation, title)
		RenderTemplate(w, r, p)
	}
}

func users(w http.ResponseWriter, r *http.Request, title string) {
	p := LoadPage(w, backendLocation, title)
	users, err := data.GetAllUsers(0)
	if err != nil {
		panic(err)
	}
	p.Info = UserlistParams{users}
	RenderTemplate(w, r, p)
}

func edit(w http.ResponseWriter, r *http.Request, title string) {
	p := LoadPage(w, backendLocation, title)
	toLoad := r.FormValue("p")
	formPage := loadPageForEdit(w, toLoad)
	p.Info = EditParams{toLoad, string(formPage.Body)}
	RenderTemplate(w, r, p)
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
	RenderTemplate(w, r, p)
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
