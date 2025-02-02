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

	// obtain the label and time associated
	game.TimeData = make(map[string]string)
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

				mainStoryData, mainStoryDataExists := game.TimeData["Main Story"]
				if (el.ChildText("h4") == "Co-Op") || (el.ChildText("h4") == "Single-Player") {
					// if main story data exists, overwrite only if new data is greater
					if (mainStoryDataExists) && (mainStoryData < el.ChildText("h5")) {
						game.TimeData["Main Story"] = el.ChildText("h5")
					} else if !mainStoryDataExists {
						// There's no Main Story data so write it
						game.TimeData["Main Story"] = el.ChildText("h5")
					}
					return
				}
				// the label is not any of the following:
				//		"All Styles"
				//		"Vs."
				//		"Co-Op"
				//		"Single-Player"
				// and we simply write the data as is
				game.TimeData[el.ChildText("h4")] = el.ChildText("h5")
			}
		})
	})

	// if the map is empty, return an empty Game Object
	// o/w grab all the rest of the data

	if len(game.TimeData) == 0 {
		game.Name = ""
		game.Url = ""
	} else {
		// obtain the game name from the page

		c.OnHTML("div.GameHeader_profile_header__q_PID", func(e *colly.HTMLElement) {
			game.Name = strings.TrimSpace(e.Text)
		})

		// when the data is acquired, log it and attach URL
		c.OnScraped(func(r *colly.Response) {
			// attach the url to the game
			game.Url = link

			fmt.Println("Data Obtained!", r.Request.URL)
		})
	}

	c.Visit(link)

	return
}
