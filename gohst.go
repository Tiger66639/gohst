package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"

	"github.com/cosban/gohst/auth"
	"github.com/cosban/gohst/data"
	"github.com/cosban/gohst/web"
	"github.com/vaughan0/go-ini"
)

var (
	conninfo string
)

func init() {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Panicln("There was an issue with the config file! ", err)
	}
	host, _ := config.Get("psql", "host")
	user, _ := config.Get("psql", "user")
	database, _ := config.Get("psql", "database")
	password, _ := config.Get("psql", "password")
	port, _ := config.Get("psql", "port")
	conninfo = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)
}

func main() {
	data.Connect(conninfo)
	l, _ := net.Listen("tcp", "127.0.0.1:8000")

	http.HandleFunc("/", web.PageHandler)
	http.HandleFunc("/static/", web.StaticHandler)
	http.HandleFunc("/txt/", web.RawHandler)
	http.HandleFunc("/dev/", web.DevHandler)
	http.HandleFunc("/backend/", web.BackendHandler)

	http.HandleFunc("/connect", auth.Connect)
	http.HandleFunc("/register/submit", auth.RegisterUser)
	http.HandleFunc("/disconnect", auth.Disconnect)

	fcgi.Serve(l, nil)
}
