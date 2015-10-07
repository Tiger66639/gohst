package web

import (
	"io/ioutil"
	"net/http"
)

type RawParams struct {
	Text string
}

// loadtxt attempts to inject a .txt file into the "raw" template.
// If the .txt file does not exist, a 404 page is displayed.
func loadRaw(w http.ResponseWriter, title string) (*Page, error) {
	body, err := ioutil.ReadFile("static/txt/" + title)
	if err != nil {
		return nil, err
	}
	p := LoadPage(w, SharedLocation, "raw")
	p.Title = title
	p.Info = RawParams{string(body)}
	return p, nil
}

// RawHandler is used for static .txt documents.
// Renders the 404 page upon error.
func RawHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/txt/"):]
	p, err := loadRaw(w, title)
	if err != nil {
		panic(err)
		RenderTemplate(w, r, OnError(w, 404))
		return
	}
	RenderTemplate(w, r, p)
}
