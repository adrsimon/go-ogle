package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./crawler.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func search(rows *sql.Rows) []string {
	var urls []string
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
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
	rows, err := db.Query("SELECT url FROM headers WHERE header LIKE ?", "%"+input+"%")
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
		println(url)
	}
}
