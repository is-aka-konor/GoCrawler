package main

import (
	"flag"
	"fmt"
	"os"

	a5e "GoCrawler/internal/models"

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
	// // Initialize the collector
	collector := colly.NewCollector()

	var infoList []a5e.SpellList

	collector.OnHTML("td.views-field-title a", func(e *colly.HTMLElement) {
		info := a5e.SpellList{
			SpellUrl: e.Attr("href"),
			Name:     e.Text,
		}
		infoList = append(infoList, info)
		fmt.Printf("Spell: %s, URL: %s\n", info.Name, info.SpellUrl)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Start the scraping process
	collector.Visit(startURL)
}
