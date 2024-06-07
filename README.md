# Go-ogle

A simple web search engine written in Golang using the [Colly](https://go-colly.org/) web crawling framework. 

## Progress

1% complete
- [x] Basic web crawling
- [x] DB saving
- [ ] Advanced parallel crawling
- [ ] Page understanding
- [ ] Ranking
- [ ] Indexing
- [ ] Web UI

## Usage

Run the web crawler using the following command:
```bash
go run cmd/crawler/crawl.go
```

To query the database, run the following command:
```bash
go run cmd/querier/query.go queryWord
```
This will report all the hits in the database that contain the word `queryWord`.