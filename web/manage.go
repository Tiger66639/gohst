package web

// type Params struct {
// 	Title string
// 	Body  []byte
// }

// /**
//  * For backend specific pages
//  */
// func BackendHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/"):]
//
// 	if !auth.IsConnected(r) {
// 		login := LoadPage(w, "login")
// 		RenderTemplate(w, login.Filename)
// 		return
// 	}
//
// 	switch title[len("backend/"):] {
// 	case "edit":
// 		p := LoadPage(w, title)
// 		toLoad := r.FormValue("p")
// 		formPage := loadPageForEdit(w, toLoad)
// 		p.Info = Params{toLoad, formPage.Body}
// 		RenderTemplate(w, p.Filename)
// 		break
// 	case "edit/submit":
// 		formURL := r.FormValue("p")
// 		formPage := loadPageForEdit(w, formURL)
// 		formPage.Body = []byte(r.FormValue("body"))
// 		formPage.Template = template.Must(template.New(formPage.Filename).Parse(string(formPage.Body)))
// 		formPage.Template = template.Must(formPage.Template.ParseFiles("templates/base.html"))
// 		pages[formPage.Filename] = formPage
// 		http.Redirect(w, r, "/backend/manage", http.StatusFound)
// 		log.Printf("SECOND\n%s", formPage.Template)
// 		break
// 	case "manage":
// 		p := LoadPage(w, title)
// 		RenderTemplate(w, p.Filename)
// 		break
// 	}
// }
//
// /**
//  * Loads a page from the cache, if page isn't in cache, loads the blank page
//  */
// func loadPageForEdit(w http.ResponseWriter, title string) *Page {
// 	if title == "base" || title == "blank" {
// 		return BlankPage(w)
// 	}
//
// 	filename := "./templates/" + title + ".html"
// 	if strings.Contains(title, "/") {
// 		title = title[strings.LastIndex(title, "/")+1:]
// 	}
//
// 	if page, ok := pages[filename]; ok {
// 		return page
// 	}
//
// 	body, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return BlankPage(w)
// 	}
//
// 	// page exists add it to the cache
// 	tmpl := template.Must(template.ParseFiles(filename, "templates/base.html"))
// 	pages[filename] = &Page{Title: title, Filename: filename, Template: tmpl, Body: body}
// 	return pages[filename]
// }
