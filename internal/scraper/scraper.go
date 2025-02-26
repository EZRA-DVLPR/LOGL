package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Game struct {
	Name, HLTBUrl, CompletionatorUrl string
	Favorite                         int
	Main, MainPlus, Comp             float32
}

// given the name of a game as a string, search HLTB, get its data and return as game struct
// if the search fails, then searches bing and gets the first hit
func SearchGameHLTB(gameName string) Game {
	log.Println("Searching HLTB for game...")

	searchRes := searchHLTB(gameName)
	if searchRes == "" {
		log.Println("Querying HLTB Failed. Retrying through Bing Search...")
		// try via bing search
		searchRes = searchBing(gameName)

		// return empty game
		if searchRes == "" {
			log.Println("No Link found. Process Aborted!")
			var emptyGame Game
			return emptyGame
		}

		log.Println("Link obtained. Web Scraping process beginning...")
		return FetchHLTB(searchRes)

	} else {
		log.Println("Link obtained. Web Scraping process beginning...")
		return FetchHLTB("https://howlongtobeat.com" + searchRes)
	}
}

// given the entire proper link for HLTB, obtain information for the game
func FetchHLTB(link string) (game Game) {
	// declare the collector object so the scraping process can begin
	c := colly.NewCollector()

	// establish connection to HLTB
	c.OnRequest(func(r *colly.Request) {
		log.Println("Connection made to HLTB")
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
					if (game.Main != 0) && (game.Main < cleanTime(el.ChildText("h5"))) {
						// There's no Main Story data so write it
						game.Main = cleanTime(el.ChildText("h5"))
					} else if game.Main == 0 {
						game.Main = cleanTime(el.ChildText("h5"))
					}
				}

				// write the data for Main Story, Main + Sides, and Completionist
				if el.ChildText("h4") == "Main Story" {
					game.Main = cleanTime(el.ChildText("h5"))
				}
				if el.ChildText("h4") == "Main + Sides" {
					game.MainPlus = cleanTime(el.ChildText("h5"))
				}
				if el.ChildText("h4") == "Completionist" {
					game.Comp = cleanTime(el.ChildText("h5"))
				}
			}
		})
		// when finished obtaining all the data, fill all empty values with "-1"

		if game.Main == 0 {
			game.Main = -1
		}
		if game.MainPlus == 0 {
			game.MainPlus = -1
		}
		if game.Comp == 0 {
			game.Comp = -1
		}
	})

	// set the game name
	c.OnHTML("div.GameHeader_profile_header__q_PID", func(e *colly.HTMLElement) {
		// remove the stuff after the <br>
		game.Name = strings.TrimSpace(e.Text)
	})

	// when the data is acquired, log it and attach URL
	c.OnScraped(func(r *colly.Response) {
		// attach the url to the game
		game.HLTBUrl = link

		// make the game not favorited
		game.Favorite = 0

		log.Println("Data Obtained!", r.Request.URL)
	})

	c.Visit(link)

	return
}

// given the name of a game as a string, search Completionator, get its data and return as game struct
func SearchGameCompletionator(gameName string) Game {
	log.Println("Searching Completionator for game...")
	searchRes := searchCompletionator(gameName)
	if searchRes == "" {
		log.Println("No Link found. Process Aborted!")
		var emptyGame Game
		return emptyGame
	} else {
		log.Println("Link obtained. Web Scraping process beginning...")
		return FetchCompletionator("https://completionator.com" + searchRes)
	}
}

func FetchCompletionator(link string) (game Game) {
	// declare the collector object so the scraping process can begin
	c := colly.NewCollector()

	// establish connection to Completionator
	c.OnRequest(func(r *colly.Request) {
		log.Println("Connection made to Completionator")
	})

	// log that there was a problem accessing the URL
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// update the Main Story, Main + Sides, and Completionist fields of the game struct
	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		e.ForEach("div.col-6", func(_ int, el *colly.HTMLElement) {
			// dont grab the data that is from the following categories:
			// 		"speed run"
			if el.ChildText("h5") != "speed run" {
				// write the data for Main Story, Main + Sides, and Completionist
				// In completionator they are saved as "core + few", "core + lots", "completionated"
				if el.ChildText("h5") == "core + few" {
					game.Main = cleanTime(el.ChildText("h3"))
				}
				if el.ChildText("h5") == "core + lots" {
					game.MainPlus = cleanTime(el.ChildText("h3"))
				}
				if el.ChildText("h5") == "completionated" {
					game.Comp = cleanTime(el.ChildText("h3"))
				}
			}
		})
		// when finished obtaining all the data, fill all empty values with "-1"

		if game.Main == 0 {
			game.Main = -1
		}
		if game.MainPlus == 0 {
			game.MainPlus = -1
		}
		if game.Comp == 0 {
			game.Comp = -1
		}
	})

	// set the game name
	c.OnHTML("h2.game-details-title", func(e *colly.HTMLElement) {
		// grab the first child of h2 tag
		game.Name = strings.TrimSpace(e.DOM.Contents().First().Text())
	})

	// when the data is acquired, log it and attach URL
	c.OnScraped(func(r *colly.Response) {
		// attach the url to the game
		game.CompletionatorUrl = link

		// make the game not favorited
		game.Favorite = 0

		log.Println("Data Obtained!", r.Request.URL)
	})

	c.Visit(link)

	return
}

func cleanTime(time string) (cleanTime float32) {
	// if no time recorded, then return -1
	if time == "--" {
		return -1
	}

	// immediately cut out all trailing data from the space including the space
	// eg. Hours
	time, _, _ = strings.Cut(time, " ")

	// if time contains "½", then add .5 to cleantime
	if strings.Contains(time, "½") {
		cleanTime = 0.5
		time, _, _ = strings.Cut(time, "½")
	}

	// convert the whole number portion to a float 32
	cleanTimeWhole, err := strconv.ParseFloat(time, 32)
	if err != nil {
		log.Println("Error converting given Clean Time to Float", err)
	}

	cleanTime = cleanTime + float32(cleanTimeWhole)
	return
}
