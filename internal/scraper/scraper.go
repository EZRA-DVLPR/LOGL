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
	Main, MainPl, Comp, All string
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

func FetchHLTB() {
	// declare the collector object
	c := colly.NewCollector(
	// TODO: Need to check how it works with different User Agents
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("HLTB connection made, r.URL")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Data Obtained!", r.Request.URL)
	})

	var games []Game

	c.OnHTML("li.GameStats_short__tSJ6I", func(e *colly.HTMLElement) {
		// make new game instance
		game := Game{}

		// scrape data
		game.Main = e.ChildText("h5")
		// game.MainPl = e.ChildAttr("img", "src")
		// game.Comp = e.ChildAttr("img", "src")
		// game.All = e.ChildAttr("img", "src")

		game.MainPl = "MainPlus"
		game.Comp = ""
		game.All = ""

		fmt.Println(game.Main)
		fmt.Println(game.MainPl)
		fmt.Println("Comp")

		// add product instance with scraped data to the list of products
		games = append(games, game)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)

		for _, game := range games {
			fmt.Println(game.Main)
			fmt.Println(game.MainPl)
		}
	})

	c.Visit("https://howlongtobeat.com/game/68151")
}
