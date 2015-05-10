package main

import (
	"net"
	"net/http"
	"net/http/fcgi"

	"github.com/cosban/gohst/auth"
	"github.com/cosban/gohst/web"
)

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:8000")
	//TODO: remove directory listing from this
	http.HandleFunc("/", web.PageHandler)
	http.HandleFunc("/static/", web.StaticHandler)
	http.HandleFunc("/txt/", web.TxtHandler)
	http.HandleFunc("/dev/", web.DevHandler)

	http.HandleFunc("/connect", auth.Connect)

	fcgi.Serve(l, nil)
}
