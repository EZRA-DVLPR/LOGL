package main

import (
	"fmt"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

func main() {
	fmt.Println("Program Start")

	// scraper.TestScraper()

	scraper.FetchHLTBRunner()

	// scraper.GPTTest()
}
