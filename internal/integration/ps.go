package integration

import (
	"context"
	"fmt"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/chromedp/chromedp"
)

func GetAllGamesPS(profile string) {
	fmt.Println("Getting games for PSN...")

	// final list holding all games from all pages
	var gamelist []string

	// get games from first page, append them into gamelist, and continue grabbing until the last page
	gamepartlist, nextpage := getAllGamesPS(profile, "1")
	for nextpage != "0" {
		gamelist = append(gamelist, gamepartlist...)
		gamepartlist, nextpage = getAllGamesPS(profile, nextpage)
	}

	// append the games retrieved from the last page
	gamelist = append(gamelist, gamepartlist...)
	for _, game := range gamelist {
		fmt.Println(game)
	}
}

func getAllGamesPS(profile string, pagenum string) (gamelist []string, nextPageNum string) {
	url := "https://psnprofiles.com/" + profile + "?ajax=1&page=" + pagenum

	// Define a custom user agent
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	// Set Chrome options for headless execution
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),    // Ensure Chrome runs headless
		chromedp.Flag("disable-gpu", true), // Disable GPU to avoid issues
		chromedp.Flag("no-sandbox", true),  // Required for some environments
		chromedp.UserAgent(userAgent),      // Set custom user agent
	)

	// Create an allocator with the defined options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Perform the search on psnprofiles
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

func isAnotherGame(pageHTML string) (isAnother bool) {
	// look for the next game name
	indexGameNameStart := strings.Index(pageHTML, `alt=\"`)

	if indexGameNameStart == -1 {
		return false
	}
	return true
}

func getAnotherGame(pageHTML string) (gameName string, nextStartIndex int) {
	// find where the game name exists
	indexGameNameStart := strings.Index(pageHTML, `alt=\"`) + 6
	indexGameNameEnd := strings.Index(pageHTML[indexGameNameStart:], `\"`) + indexGameNameStart

	// clean the unicode leftover
	gameName, err := strconv.Unquote(`"` + pageHTML[indexGameNameStart:indexGameNameEnd] + `"`)
	if err != nil {
		fmt.Println("error cleaning unicode text:", err)
		return pageHTML[indexGameNameStart:indexGameNameEnd], indexGameNameEnd
	}
	// clean the & symbol that is leftover
	gameName = strings.Replace(gameName, "&amp;", "&", -1)

	return gameName, indexGameNameEnd
}

func getNextPage(pageHTML string) (nextPage string) {
	// find NextPage
	indexNP := strings.Index(pageHTML, "nextPage = ") + 11
	// find the first characters after the index from nextpage
	indexEndNP := strings.Index(pageHTML[indexNP:], "\\r") + indexNP

	return pageHTML[indexNP:indexEndNP]
}

func decodeHTMLEnts(text string) string {
	// unescape common html entities
	text = html.UnescapeString(text)

	// matches things like `&#039`
	regex := regexp.MustCompile(`&#(\d+);`)

	// replaces the numerical character subset (unicode) with the characters that they represent
	// eg. #039 becomes '
	return regex.ReplaceAllStringFunc(text, func(entity string) string {
		// get the numbers
		numStr := strings.Trim(entity, "&#;")
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return entity
		}
		return string(rune(num))
	})
}
