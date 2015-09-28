package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Page struct {
	PageId   int
	Title    string
	Disabled bool
}

func GetAllPages() ([]*Page, error) {
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM PAGES").Scan(&count)
	rows, err := db.Query("SELECT * FROM PAGES")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var pages = make([]*Page, count)
	i := 0
	for rows.Next() {
		var pageid int
		var title string
		var disabled bool
		err = rows.Scan(&pageid, &title, &disabled)
		if err != nil {
			return nil, err
		}
		pages[i] = &Page{pageid, title, disabled}
		i++
	}
	return pages, err
}

func GetPageFromTitle(title string) (*Page, error) {
	var page *Page
	err := db.QueryRow("SELECT * FROM PAGES WHERE TITLE=$1", title).Scan(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func GetPageFromPageId(pageid int) (*Page, error) {
	var page *Page
	err := db.QueryRow("SELECT * FROM PAGES WHERE pageid=$1", pageid).Scan(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func AddNewPage(page *Page) sql.Result {
	result, err := db.Exec("INSERT INTO PAGES (title, disabled) VALUES ($1,$2)", page.Title, page.Disabled)
	if err != nil {
		panic(err)
	}
	return result
}

func RemovePage(page *Page) sql.Result {
	result, err := db.Exec("DELETE FROM PAGES WHERE pageid=$1", page.PageId)
	if err != nil {
		panic(err)
	}
	return result
}
