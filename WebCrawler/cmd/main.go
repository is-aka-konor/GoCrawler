package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"WebCrawler/internal/models"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Define command-line flags
	domainFlag := flag.String("domain", "", "The allowed domain for scraping")
	startURLFlag := flag.String("startURL", "", "The URL where the scraping starts")

	// Parse the command-line flags
	flag.Parse()

	// Check if the flags were provided, otherwise use environment variables
	domain := *domainFlag
	if domain == "" {
		domain = os.Getenv("DOMAIN")
	}
	startURL := *startURLFlag
	if startURL == "" {
		startURL = os.Getenv("START_URL")
	}

	// Check if the domain and startURL are provided
	if domain == "" || startURL == "" {
		fmt.Println("Both domain and startURL arguments are required, either as flags or environment variables")
		return
	}
	settings := models.ParserSettings{
		StartPoint: 0,
		EndPoint:   8,
		QueryParam: "?combine=&field_spell_ritual_value=All&page=",
		BaseURL:    startURL,
	}
	// // Initialize the collector
	collector := colly.NewCollector()
	// Add a random delay to the requests
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*a5e.tool*",
		Parallelism: 4,
		RandomDelay: 10 * time.Second,
	})

	infoList := make([]models.SpellList, 0, 450)

	collector.OnHTML("td.views-field-title a", func(e *colly.HTMLElement) {
		info := models.SpellList{
			SpellUrl: e.Attr("href"),
			Name:     e.Text,
		}
		infoList = append(infoList, info)
		fmt.Printf("Spell: %s, URL: %s\n", info.Name, info.SpellUrl)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start the scraping process
	for i := settings.StartPoint; i <= settings.EndPoint; i++ {
		collector.Visit(settings.BaseURL + settings.QueryParam + fmt.Sprint(i))
	}

	dataCollector := collector.Clone()

	dataCollector.OnHTML("h1.page-header", func(e *colly.HTMLElement) {
		spell := models.Spell{
			Name: e.Text,
		}
		fmt.Println("Spell:", spell.Name)
	})

	dataCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Data Collector Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	dataCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Data Collector is visiting", r.URL)
	})

	dataCollector.OnResponse(func(r *colly.Response) {
		error := r.Save(fmt.Sprintf("../FileCrawler/html/%s.html", strings.Replace(r.FileName(), ".unknown", "", -1)))
		if error != nil {
			fmt.Println("Error saving file: ", error)
		}
	})

	for _, info := range infoList {
		dataCollector.Visit(fmt.Sprintf("%s%s", domain, info.SpellUrl))
	}

	content, err := json.MarshalIndent(infoList, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return
	}

	// Write the JSON data to a file
	err = os.WriteFile("spells.json", content, 0755)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}

}
