package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type jobDetail struct {
	title    string
	location string
	salary   string
	summary  string
	dataJK   string
}

// Scrape jobs from an URL
func Scrape(baseURL string) {
	totalPages := getTotalPages(baseURL)
	var jobs []jobDetail

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, getPage(i, baseURL)...)

	}

	fmt.Println(jobs)
}

// Extract job from a card
func extractJob(card *goquery.Selection) jobDetail {

	title := cleanString(card.Find(".title").Text())
	location := cleanString(card.Find(".location").Text())
	salary := cleanString(card.Find(".salary").Text())
	summary := cleanString(card.Find(".summary").Text())
	dataJK, _ := card.Attr("data-jk")

	return jobDetail{title, location, salary, summary, dataJK}
}

// Get jobs on a page
func getPage(i int, baseURL string) []jobDetail {

	var jobs []jobDetail

	// GET request to the URL
	url := baseURL + "&start=" + strconv.Itoa(i*50)
	res, err := http.Get(url)
	checkError(err)
	checkCode(res.StatusCode)

	// Close the body when func's done
	defer res.Body.Close()

	// Get from document from the response body
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	// Find all the job cards
	cards := doc.Find(".jobsearch-SerpJobCard")

	// Extract information from each card
	cards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

// Get total pages
func getTotalPages(url string) int {

	// GET request to the URL
	res, err := http.Get(url)
	checkError(err)
	checkCode(res.StatusCode)

	// Close the body when func's done
	defer res.Body.Close()

	// Get from document from the response body
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	// Find the pagination
	pagination := doc.Find(".pagination")

	// Get tne count of pages
	totalPages := 0
	pagination.Each(func(i int, pageList *goquery.Selection) {
		totalPages = pageList.Find("a").Length()
	})

	return totalPages
}

// Check Error
func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Check Status Code
func checkCode(statusCode int) {
	if statusCode != 200 {
		log.Fatalln(statusCode)
	}
}

func cleanString(str string) string {
	return strings.TrimSpace(str)
}
