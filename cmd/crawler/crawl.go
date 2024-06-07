package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

var queue []string
var cooldown = map[string]int{}
var c *colly.Collector

func getDomain(url string) string {
	parts := strings.Split(url, "/")
	domain := parts[2]

	n := len(strings.Split(domain, "."))
	if n == 2 {
		return domain
	} else {
		return strings.Split(domain, ".")[n-2] + "." + strings.Split(domain, ".")[n-1]
	}
}

func setupDatabase() *sql.DB {
	println("Setting up database...")
	db, err := sql.Open("sqlite3", "./crawler.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS headers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL UNIQUE,
        header TEXT NOT NULL
    );`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	println("Database setup done !")
	return db
}

func insertHeader(db *sql.DB, url, header string) {
	stmt, err := db.Prepare("INSERT INTO headers(url, header) VALUES(?, ?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		log.Fatal(err)
	}
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(stmt)

	_, err = stmt.Exec(url, header)
	if err != nil {
		log.Fatal(err)
	}
}

func saveQueue(queue []string) {
	f, err := os.Create("queue.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	for _, url := range queue {
		_, err = f.WriteString(url + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func retrieveQueue() []string {
	println("Retrieving queue...")

	f, err := os.Open("queue.txt")
	if err != nil {
		println("No queue found, starting from scratch...")
		return []string{"https://www.lemonde.fr"}
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	var retrieved []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		retrieved = append(retrieved, scanner.Text())
	}

	println("Queue retrieved !")
	return retrieved
}

func main() {
	db := setupDatabase()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	c = colly.NewCollector(
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if len(link) > 0 && strings.HasPrefix(link, "https://") {
			queue = append(queue, link)
		}
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		header := e.Text
		url := e.Request.URL.String()

		insertHeader(db, url, header)
	})

	queue = retrieveQueue()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for len(queue) > 0 {
			url := queue[0]
			queue = queue[1:]
			if cooldown[getDomain(url)] > 0 {
				continue
			}
			fmt.Println("queue length: ", len(queue))
			fmt.Println("visiting domain: ", getDomain(url))
			err := c.Visit(url)
			if err != nil {
				fmt.Println("Error visiting", url, ":", err)
			}

			cooldown[getDomain(url)] = 150

			for domain, time := range cooldown {
				cooldown[domain] = time - 1
				if time == 0 {
					delete(cooldown, domain)
				}
			}

			c.Wait()
		}
		fmt.Println("Done")
	}()

	<-sigChan
	fmt.Println("Interrupting...")
	saveQueue(queue)
	fmt.Println("Queue saved")
	fmt.Println("Exiting...")
}
