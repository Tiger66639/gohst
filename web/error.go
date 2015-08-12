package web

import (
	"log"
	"net/http"
)

// PageError is the additional information needed to display the error page
type PageError struct {
	Type        int
	Description string
}

var definitions = map[int]string{
	404: "That page could not be found",
	401: "Login to view this page",
	403: "Unauthorized access",
}

// OnError is called whenever an error occurs that the server is able to handle
func OnError(w http.ResponseWriter, status int) *Page {
	log.Printf("%d", status)
	w.WriteHeader(status)
	p := LoadPage(w, SharedLocation, "error")
	p.Info = PageError{status, definitions[status]}
	return p
}
