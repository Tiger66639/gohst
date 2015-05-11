package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cosban/gohst/auth"
)

/**
 * represents all of the information needed to render a page
 */
type Page struct {
	// The title of the page and its filename
	Title, Filename string
	// The html template
	Template *template.Template
	// A byte array of the body
	Body []byte
	// strings needed in order for oauth to work
	// TODO: understand this better
	ClientID, State, Scope string
	// AccessToken needed for cookies
	AccessToken string
}

/**
 * A map of pages keyed by their filename.
 * This caches templates to reduce CPU load
 */
var pages = make(map[string]*Page)

/**
 * The handler used for loading static content.
 * Since we should not be internally requesting non-existant content,
 * a 404 page is rendered upon failure to directory prevent probing.
 */
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}
	PageHandler(w, r)
}

/**
 * The standard page handler used for most cases
 */
func PageHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	if len(title) == 0 {
		title = "home"
	}
	p := LoadPage(w, title)
	RenderTemplate(w, p.Filename)
}

/**
 * Handler used when authentication is required to load a page
 */
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if auth.IsConnected(r) {
		log.Printf("Found a valid receipt")
		PageHandler(w, r)
		return
	}
	log.Printf("No valid receipt found")
	login := LoadPage(w, "login")
	RenderTemplate(w, login.Filename)
}

/**
 * In the event that /dev is placed within the url, this refreshes the cache
 * with the most updated html template.
 */
func DevHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/dev/"):]
	if auth.IsConnected(r) {
		delete(pages, "./templates/"+title+".html")
	}
	p := LoadPage(w, title)
	RenderTemplate(w, p.Filename)
}

/**
 * Checks if the file is within the cache, returning from it if s
 * If not, it checks whether the file exists and saves it to the cache
 * while also returning it.
 * If the file does not exist, a 404 error is thrown and the 404 page is
 * rendered.
 */
func LoadPage(w http.ResponseWriter, title string) *Page {
	filename := "./templates/" + title + ".html"
	if strings.Contains(title, "page/") {
		title = title[len("page/"):]
	}
	// if the page is inside the cache, just load it
	if page, ok := pages[filename]; ok {
		return page
	}
	// the page is not inside the cache so see if it exists
	body, err := ioutil.ReadFile(filename)
	// throw a 404 error
	if err != nil {
		return OnError(w)
	}
	// page exists add it to the cache
	tmpl := template.Must(template.ParseFiles(filename, "templates/base.html"))
	pages[filename] = &Page{Title: title, Filename: filename, Template: tmpl, Body: body}
	return pages[filename]
}

/**
 * The handler used for static .txt documents.
 * Renders the 404 page upon error.
 */
func TxtHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/txt/"):]
	p, err := Loadtxt(w, title)
	if err != nil {
		RenderTemplate(w, OnError(w).Filename)
		return
	}
	fmt.Fprintf(w, "%s", p.Body)
}

/**
 * Loads the content from static .txt files to be used with the handler
 */
func Loadtxt(w http.ResponseWriter, title string) (*Page, error) {
	body, err := ioutil.ReadFile("static/txt/" + title)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: []byte("<html><body><pre>" + string(body) + "</pre></body></html>")}, nil
}

/**
 * Renders the page template using all of the data requested
 */
func RenderTemplate(w http.ResponseWriter, name string) {
	pages[name].Template.ExecuteTemplate(w, "base", pages[name])
}

/**
 * Called on any content error, loads and renders the 404 page
 */
func OnError(w http.ResponseWriter) *Page {
	w.WriteHeader(404)
	return LoadPage(w, "404")
}
