package main

import (
	"FileCrawler/internal/spells"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	// define command-line flags
	folderFlag := flag.String("folder", "", "The folder to crawl")
	// parse the command-line flags
	flag.Parse()

	// check if the flags were provided, otherwise use environment variables
	folder := *folderFlag
	if folder == "" {
		folder = os.Getenv("FOLDER")
	}

	// check if the folder is provided
	if folder == "" {
		fmt.Println("Folder argument is required, either as a flag or environment variable")
		return
	}

	// validate the folder path exists
	dir, err := filepath.Abs(filepath.Dir(folder))
	if err != nil {
		fmt.Println("Invalid folder path:", folder)
		panic(err)
	}

	pages, err := getFileList(dir)
	if err != nil {
		fmt.Println("Error getting file list:", err)
		return
	}

	// initialize the collector
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	collector := colly.NewCollector()
	collector.WithTransport(transport)

	spellList := make([]*spells.Spell, 0, len(pages))
	collector.OnHTML(".page-content", func(e *colly.HTMLElement) {
		var spell = spells.Spell{}

		spell.Name = e.ChildText("h1.page-header")
		spell.Level, err = getSpellLevel(e)
		if err != nil {
			fmt.Println("Error getting spell level:", err)
		}
		// Print the whole spell struct to the console and append it to the spellList
		fmt.Printf("Spell: %+v\n", spell)
		spellList = append(spellList, &spell)
	})

	// crawl the folder
	for _, page := range pages {
		// Convert Windows file path to a valid file URL
		fileURL := "file://" + strings.Replace(page, "\\", "/", -1)
		fmt.Printf("Crawling file: %s \n", fileURL)
		err := collector.Visit(fileURL)
		if err != nil {
			fmt.Println("Error crawling file:", err)
		}
		collector.Wait()
	}

}

func getFileList(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileList := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			subFiles, err := getFileList(subDir)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, subFiles...)
		} else {
			fileList = append(fileList, filepath.Join(dir, file.Name()))
		}
	}

	return fileList, nil
}

// Helper functions to pull data from the HTML pages
func getSpellLevel(e *colly.HTMLElement) (int, error) {
	level := -1
	levelTxt := e.ChildText(".field--name-field-spell-level a")
	if levelTxt != "" {
		levelTxt = strings.TrimSpace(levelTxt)
		if strings.EqualFold(levelTxt, "Cantrip") {
			level = 0
			return level, nil
		}
		// Remove the last two characters from levelText
		if len(levelTxt) > 2 {
			levelTxt = levelTxt[:len(levelTxt)-2]
		}

		// Attempt to parse the remaining string into an integer
		level, err := strconv.Atoi(strings.TrimSpace(levelTxt))
		if err != nil {
			// Handle the error, for example, by returning a default value or forwarding the error
			return 0, err
		}

		return level, nil
	}
	err := fmt.Errorf("no level found")
	return level, err
}
