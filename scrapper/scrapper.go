package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
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

var fileName = "jobs.csv"

// Scrape jobs from an URL
func Scrape(baseURL string) {
	totalPages := getTotalPages(baseURL)
	var jobs []jobDetail

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, getPage(i, baseURL)...)

	}

	writeJobs(jobs)

	fmt.Println("Done")
}

func extractDetail(i int, job jobDetail) []string {
	url := "https://www.indeed.com/viewjob?jk="
	id := strconv.Itoa(i + 1)
	link := url + job.dataJK
	return []string{id, job.title, job.location, job.salary, job.summary, link}
}

func writeJobs(jobs []jobDetail) {

	file, err := os.Create(fileName)
	checkError(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	rows := [][]string{}

	rows = append(rows, []string{"ID", "TITLE", "LOCATION", "SALARY", "SUMMARY", "LINK"})

	for i, job := range jobs {
		rows = append(rows, extractDetail(i, job))
	}

	wErr := w.WriteAll(rows)
	checkError(wErr)
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
