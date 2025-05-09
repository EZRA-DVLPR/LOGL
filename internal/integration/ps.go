package integration

import (
	"context"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/model"
	"github.com/chromedp/chromedp"
)

func GetAllGamesPS(profile string) {
	log.Println("Getting games for PSN")

	// final list holding all games from all pages
	var gameList []string

	// get games from first page, append them into gamelist, and continue grabbing until the last page
	gamepartlist, nextpage := getAllGamesPS(profile, "1")
	log.Println("Obtained all game titles from page 1")
	for nextpage != "0" {
		log.Println("Obtaining list of game titles from page: " + nextpage)
		gameList = append(gameList, gamepartlist...)              // unpack and append each elt from part to gameList
		gamepartlist, nextpage = getAllGamesPS(profile, nextpage) // get next page
		log.Println("Obtained all game titles from page: ", nextpage)
	}

	// append the games retrieved from the last page retrieved
	gameList = append(gameList, gamepartlist...)
	log.Println("Obtained all game titles for profile:", profile)
	model.SetMaxProcesses(len(gameList))
	for _, game := range gameList {
		dbhandler.SearchAddToDB(game)
	}
	log.Println("Finished adding game data from PSN for profile:", profile)
}

func getAllGamesPS(profile string, pagenum string) (gamelist []string, nextPageNum string) {
	url := "https://psnprofiles.com/" + profile + "?ajax=1&page=" + pagenum

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent(userAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pageHTML string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &pageHTML),
	)
	if err != nil {
		log.Fatal(err)
	}

	// clean up the html unicode stuff into their respective characters
	pageHTML = decodeHTMLEnts(pageHTML)

	// while loop grabs games so long as there is another game
	for isAnotherGame(pageHTML) {
		nextGameName, newStartIndex := getAnotherGame(pageHTML)
		gamelist = append(gamelist, nextGameName)
		pageHTML = pageHTML[newStartIndex:]
	}

	return gamelist, getNextPage(pageHTML)
}

// look for the next game. if exists then true. o/w false
func isAnotherGame(pageHTML string) (isAnother bool) {
	indexGameNameStart := strings.Index(pageHTML, `alt=\"`)

	if indexGameNameStart == -1 {
		return false
	}
	return true
}

// will get the game name for the next game, returns the name and index after the game name
func getAnotherGame(pageHTML string) (gameName string, nextStartIndex int) {
	// find where the game name exists
	indexGameNameStart := strings.Index(pageHTML, `alt=\"`) + 6
	indexGameNameEnd := strings.Index(pageHTML[indexGameNameStart:], `\"`) + indexGameNameStart

	// clean the unicode leftover
	gameName, err := strconv.Unquote(`"` + pageHTML[indexGameNameStart:indexGameNameEnd] + `"`)
	if err != nil {
		log.Println("error cleaning unicode text:", err)
		return pageHTML[indexGameNameStart:indexGameNameEnd], indexGameNameEnd
	}
	// clean the & symbol that is leftover
	gameName = strings.Replace(gameName, "&amp;", "&", -1)

	return gameName, indexGameNameEnd
}

// find NextPage and returns the value for it as a string
func getNextPage(pageHTML string) (nextPage string) {
	indexNP := strings.Index(pageHTML, "nextPage = ") + 11
	indexEndNP := strings.Index(pageHTML[indexNP:], "\\r") + indexNP

	return pageHTML[indexNP:indexEndNP]
}

// cleans text to remove wonky representations of characters (\u###), (&amp;####;), etc.
func decodeHTMLEnts(text string) string {
	// unescape common html entities
	text = html.UnescapeString(text)

	// matches things like `&#039`
	regex := regexp.MustCompile(`&#(\d+);`)

	// replaces the numerical character subset (unicode) with the characters that they represent
	// eg. #039 becomes '
	return regex.ReplaceAllStringFunc(text, func(entity string) string {
		// get the numbers then return the character they represent
		numStr := strings.Trim(entity, "&#;")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return entity
		}
		return string(rune(num))
	})
}
