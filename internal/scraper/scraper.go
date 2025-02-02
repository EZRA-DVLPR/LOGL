package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Game struct {
	Name, Url string
	TimeData  map[string]string
	// Labels    []string
	// Lengths   []string
}

var games []Game

// given the name of a game as a string, search HLTB, then get its data
func SearchGame(gameName string) {
	// TODO: Check if the game exists in the current database
	// if not then add a new entry to the database

	FetchHLTBRunner(SearchHLTB(gameName))

	fmt.Println("Game data acquired for ", gameName)
}

// given a link to a particular game, fetch the data for it
// eg. /game/42069
func FetchHLTBRunner(gameLink string) {
	games = append(games, FetchHLTB("https://howlongtobeat.com/"+gameLink))

	for index, game := range games {
		fmt.Printf("Game %d: Name: %s URL:%s \n", index+1, game.Name, game.Url)

		for label, length := range game.TimeData {
			fmt.Println(label, length)
		}

		// fmt.Printf("Labels %s: \n", game.Labels)
		// fmt.Printf("Lengths %s: \n", game.Lengths)
		fmt.Println()
	}
}

// given the entire proper link for HLTB, obtain labels and lengths for the game
func FetchHLTB(link string) (game Game) {
	// declare the collector object so the scraping process can begin
	c := colly.NewCollector()

	// establish connection to HLTB
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Connection made to HLTB")
	})

	// log that there was a problem accessing the URL
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// obtain the game name from the page
	c.OnHTML("div.GameHeader_profile_header__q_PID", func(e *colly.HTMLElement) {
		game.Name = strings.TrimSpace(e.Text)
	})

	// obtain the label and time associateD
	game.TimeData = make(map[string]string)
	c.OnHTML("div.GameStats_game_times__KHrRY", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			// TODO: Make it a setting such that if the setting is set, it will accept `All Styles` and add the data to the Game
			if el.ChildText("h4") != "All Styles" {
				game.TimeData[el.ChildText("h4")] = el.ChildText("h5")
			}

			// game.Labels = append(game.Labels, el.ChildText("h4"))   // get the label for the time eg. "Main Story"
			// game.Lengths = append(game.Lengths, el.ChildText("h5")) // get the time eg. "4 Hours"
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
