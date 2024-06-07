package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	url    string
	header string
}

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./crawler.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func search(rows *sql.Rows) []Data {
	var urls []Data
	for rows.Next() {
		var url string
		var header string

		err := rows.Scan(&url, &header)
		if err != nil {
			log.Fatal(err)
		}

		urls = append(urls, Data{
			url:    url,
			header: strings.TrimSpace(header),
		})
	}
	return urls
}

func main() {
	db := setupDatabase()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	input := os.Args[1]
	rows, err := db.Query("SELECT url, header FROM headers WHERE header LIKE ?", "%"+input+"%")
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	urls := search(rows)
	for _, url := range urls {
		println(url.url + " : " + url.header)
	}
}
