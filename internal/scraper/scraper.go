// https://www.zenrows.com/blog/web-scraping-golang#build-first-golang-scraper
package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

type Game struct {
	Name, Url string
	Labels    []string
	Lengths   []string
}

func TestScraper() {
	// declare the collector object
	c := colly.NewCollector(
	// colly.AllowedDomains("www.scrapingcourse.com"),
	)

	var products []Product

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting, r.URL")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		// make new product instance
		product := Product{}

		// scrape data
		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		// add product instance with scraped data to the list of products
		products = append(products, product)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)

		for _, product := range products {
			fmt.Println(product.Name)
		}
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce")
}

func FetchHLTBRunner() {
	var games []Game

	games = append(games, FetchHLTB("https://howlongtobeat.com/game/68151"))
	games = append(games, FetchHLTB("https://howlongtobeat.com/game/64753"))
	games = append(games, FetchHLTB("https://howlongtobeat.com/game/80199"))
	games = append(games, FetchHLTB("https://howlongtobeat.com/game/4249"))
	games = append(games, FetchHLTB("https://howlongtobeat.com/game/147712"))

	for index, game := range games {
		fmt.Printf("Game %d: Name: %s URL:%s \n", index+1, game.Name, game.Url)
		fmt.Printf("Labels %s: \n", game.Labels)
		fmt.Printf("Lengths %s: \n", game.Lengths)
		fmt.Println()
	}
}

func FetchHLTB(link string) (game Game) {
	// declare the collector object so the scraping process can begin
	c := colly.NewCollector(
	// TODO: Need to check how it works with different User Agents
	)

	// establish connection to HLTB
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Connection made to HLTB")
	})

	// log that there was a problem accessing the URL
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// obtain the game name from the link
	c.OnHTML("div.GameHeader_profile_header__q_PID", func(e *colly.HTMLElement) {
		game.Name = e.Text
	})

	// obtain the label and time associated
	c.OnHTML("div.GameStats_game_times__KHrRY", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			game.Labels = append(game.Labels, el.ChildText("h4"))   // get the label for the time eg. "Main Story"
			game.Lengths = append(game.Lengths, el.ChildText("h5")) // get the time eg. "4 Hours"
		})
	})

	// when the deed is done, log it and attach URL
	c.OnScraped(func(r *colly.Response) {
		// attach the url to the game
		game.Url = link

		fmt.Println("Data Obtained!", r.Request.URL)
	})

	c.Visit(link)

	return
}
