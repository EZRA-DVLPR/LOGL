package main

import (
	"fmt"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

func main() {
	fmt.Println("Program Start")

	// scraper.TestScraper()

	// scraper.FetchHLTBRunner()

	fmt.Println(scraper.SearchHLTB("cookie"))
	fmt.Println(scraper.SearchHLTB("nonsensegame that wont result in anything"))
}
