package web

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/cosban/gohst/auth"
)

// Page is a struct of all of the data needed to serve and display a
// publicly accessable page.
// Each Page is cached within the map pages after first load in order to reduce
// the time it takes to load it in the future.
type Page struct {
	// The title of the page and its filename
	Title string
	// The html template
	Template *template.Template
	// A byte array of the body
	Body []byte
	// Info is an interface which allows for additional data to be fed into pages
	Info interface{}
	// Whether the page is enabled or not
	Disabled, LoggedIn bool
}

// SharedLocation is where shared templates are (errors, base, etc)
const SharedLocation string = "templates/shared/"

// PageLocation is where public templates are
const PageLocation string = "templates/public/"

// pages is a map of pages keyed by the location of the page.
// The "home" page, would be keyed as ""/templates/public/home" for instance.
var pages = make(map[string]*Page)

// RenderTemplate executes templates which have been stored within the pages map
func RenderTemplate(w http.ResponseWriter, page *Page) {
	page.Template.ExecuteTemplate(w, "base", page)
}

// AuthHandler is used to verify that a client is logged in.
// If they are not, they are instead redirected to the login page.
// TODO: after the login page, it should direct the user to the page they
// requested originally
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if auth.IsConnected(r) {
		PageHandler(w, r)
		return
	}
	RenderTemplate(w, LoadPage(w, PageLocation, "login"))
}

// StaticHandler provides a way for static content to be served to clients.
// Static files are things like public images, javascript, or css.
// If a file is not found, a 404 should be served to prevent indexing.
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}
	PageHandler(w, r)
}

// PageHandler is the standard
func PageHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	if len(title) == 0 {
		title = "home"
	}
	var p *Page
	if title == "base" {
		p = OnError(w, 404)
	} else {
		p = LoadPage(w, PageLocation, title)
	}
	p.LoggedIn = auth.IsConnected(r)
	RenderTemplate(w, p)
}

func baseTemplate() *template.Template {
	return template.Must(template.ParseFiles(SharedLocation + "base.html"))
}

// LoadPage checks if the file is within the cache, returning from it if s
// If not, it checks whether the file exists and saves it to the cache
// while also returning it.
// If the file does not exist, a 404 error is thrown and the 404 page is
// rendered.
func LoadPage(w http.ResponseWriter, location string, title string) *Page {
	filename := location + title + ".html"
	if strings.Contains(title, "/") {
		title = title[strings.LastIndex(title, "/")+1:]
	}

	// if the page is inside the cache, just load it
	if page, ok := pages[filename]; ok {
		if page.Disabled {
			return OnError(w, 404)
		}
		return page
	}

	// the page is not inside the cache so see if it exists
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return OnError(w, 404)
	}

	// page exists add it to the cache
	tmpl := template.Must(template.ParseFiles(filename, SharedLocation+"base.html"))
	pages[filename] = &Page{Title: title, Template: tmpl, Body: body}
	return pages[filename]
}

// BlankPage is used to load just the shared "mega" template
func BlankPage(w http.ResponseWriter) *Page {
	if page, ok := pages["blank"]; ok {
		LoadPage(w, PageLocation, "blank")
		return page
	}
	tmpl := template.Must(template.ParseFiles(SharedLocation + "base.html"))
	pages["blank"] = &Page{Title: "", Template: tmpl, Body: []byte("")}
	return pages["blank"]
}

// DevHandler is used to refresh a specific page within the page cache
// A client must be authenticated before they are able to use this. If they are
// not authenticated they will simply be redirected to the cached version of the
// page
func DevHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/dev/"):]
	if auth.IsConnected(r) {
		delete(pages, PageLocation+title+".html")
		p := LoadPage(w, PageLocation, title)
		RenderTemplate(w, p)
	} else {
		title = "/" + title
		http.Redirect(w, r, title, http.StatusFound)
	}
}
