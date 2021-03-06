package web

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/cosban/gohst/auth"
	"github.com/cosban/gohst/data"
)

type UserlistParams struct {
	Users []*data.User
}

type EditParams struct {
	Title, Body string
}

type ManageParams struct {
	Pages []*data.Page
}

type ProfileParams struct {
	User *data.User
}

func FeedUsers(w http.ResponseWriter, r *http.Request) interface{} {
	return UserlistParams{data.GetAllUsers(0)}
}

func FeedEdit(w http.ResponseWriter, r *http.Request) interface{} {
	toLoad := r.FormValue("p")
	formPage := loadPageForEdit(w, toLoad)
	return EditParams{toLoad, string(formPage.Body)}
}

func FeedManage(w http.ResponseWriter, r *http.Request) interface{} {
	pages, err := data.GetAllPages()
	if err != nil {
		panic(err)
	}
	return ManageParams{pages}
}

func FeedProfile(w http.ResponseWriter, r *http.Request) interface{} {
	var user *data.User
	var err error
	var i int
	i, err = strconv.Atoi(r.FormValue("u"))
	if err != nil {
		i = 0
	}
	user = data.GetUserFromId(i)
	if user == nil {
		user = auth.GetConnectedUser(r)
	}
	return ProfileParams{user}
}

// Loads a page from the cache, if page isn't in cache, loads the blank page
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
