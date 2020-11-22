package main

import (
	"github.com/b0mink/Job-Scrapper/scrapper"
)


func main() {
	baseURL := "https://www.indeed.com/jobs?q=python&limit=50"
	scrapper.Scrape(baseURL)
}