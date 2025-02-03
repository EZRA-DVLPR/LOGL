package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Game struct {
	Name, Url, Main, MainPlus, Comp string
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
		fmt.Printf("Game %d: Name: %s URL: %s\n", index+1, game.Name, game.Url)
		fmt.Println("Main Story:\t", game.Main)
		fmt.Println("Main + Sides:\t", game.MainPlus)
		fmt.Println("Completionist:\t", game.Comp)
		fmt.Println()
	}
}

// given the entire proper link for HLTB, obtain information for the game
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

	// update the Main Story, Main + Sides, and Completionist fields of the game struct
	c.OnHTML("div.GameStats_game_times__KHrRY", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			// dont grab the data that is from the following categories:
			// 		"All Styles"
			// 		"Vs."
			if (el.ChildText("h4") != "All Styles") && (el.ChildText("h4") != "Vs.") {
				// if the current label is "Co-Op" or "Single-Player"
				// check if there is a value for "Main Story"
				// 			if true: compare the values and take the higher
				// 			else: make "Main Story" data
				// else write the the data as is

				if (el.ChildText("h4") == "Co-Op") || (el.ChildText("h4") == "Single-Player") {
					// if main story data exists, overwrite only if new data is greater

					if (game.Main != "") && (game.Main < el.ChildText("h5")) {
						// There's no Main Story data so write it
						game.Main = el.ChildText("h5")
					} else if game.Main == "" {
						game.Main = el.ChildText("h5")
					}
				}

				// write the data for Main Story, Main + Sides, and Completionist
				if el.ChildText("h4") == "Main Story" {
					game.Main = el.ChildText("h5")
				}
				if el.ChildText("h4") == "Main + Sides" {
					game.MainPlus = el.ChildText("h5")
				}
				if el.ChildText("h4") == "Completionist" {
					game.Comp = el.ChildText("h5")
				}
			}
		})
		// when finished obtaining all the data, fill all empty values with "--"

		if game.Main == "" {
			game.Main = "--"
		}
		if game.MainPlus == "" {
			game.MainPlus = "--"
		}
		if game.Comp == "" {
			game.Comp = "--"
		}
	})

	// set the game name
	c.OnHTML("div.GameHeader_profile_header__q_PID", func(e *colly.HTMLElement) {
		game.Name = strings.TrimSpace(e.Text)
	})

	// when the data is acquired, log it and attach URL
	c.OnScraped(func(r *colly.Response) {
		// attach the url to the game
		game.Url = link

		fmt.Println("Data Obtained!", r.Request.URL)
	})

	c.Visit(link)

	return
}
