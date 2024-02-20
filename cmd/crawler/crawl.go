package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
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

func main() {
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

	queue = append(queue, "https://fr.wikipedia.org/")

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
}
