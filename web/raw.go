package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// loadtxt attempts to inject a .txt file into the "raw" template.
// If the .txt file does not exist, a 404 page is displayed.
func loadRaw(w http.ResponseWriter, title string) (*Page, error) {
	body, err := ioutil.ReadFile("static/txt/" + title)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: []byte("<html><body><pre>" + string(body) + "</pre></body></html>")}, nil
}

// RawHandler is used for static .txt documents.
// Renders the 404 page upon error.
// TODO: RenderTemplate needs to be cloned into a RenderRawTemplate
func RawHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/txt/"):]
	p, err := loadRaw(w, title)
	if err != nil {
		RenderTemplate(w, OnError(w))
		return
	}
	fmt.Fprintf(w, "%s", p.Body)
}
