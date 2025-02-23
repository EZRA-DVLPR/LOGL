package scraper

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// given a name for a game, returns the link for the game
// eg. /game/68151
func searchHLTB(query string) (gameLink string) {
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

	// 3 second timeout in the event there is no game found from search
	ctx, cancel := context.WithTimeout(allocCtx, 5*time.Second)
	defer cancel()

	// Create a new Chrome context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Perform the search on HLTB
	var pageHTML string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.howlongtobeat.com/?q="+query),
		chromedp.WaitVisible(`.GameCard_inside_blur__cP8_l`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &pageHTML),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatal("No such game found within the timeout of 3 seconds!")
		} else {
			log.Fatal(err)
		}
	}

	// extract the link to the first game in the list
	return extractLinkHLTB(pageHTML)
}

func extractLinkHLTB(pageHTML string) (gameLink string) {
	// find the location where the first item from the search list is
	firstindex := strings.Index(pageHTML, "GameCard_inside_blur__cP8_l")

	// if there is no such item, then return empty string
	if firstindex == -1 {
		return ""
	}

	// navigate to where `a href=` is in the next part of the string
	ahrefindex := strings.Index(pageHTML[firstindex:], "a href=")

	// cut off all trailing info
	suffindex := strings.Index(pageHTML[firstindex+ahrefindex+7:], ">")

	// grab the first link for game data that we find
	// we add 8 to cut off the `a href="` from the search
	// we also add 7 for the same reason as above, but also subtract 1
	// so that it can remove the `"`
	return pageHTML[firstindex+ahrefindex+8 : firstindex+ahrefindex+7+suffindex-1]
}

// given a name for a game, returns the link for the game
// eg. /Game/Details/3441
func searchCompletionator(query string) (gameLink string) {
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent(userAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := context.WithTimeout(allocCtx, 3*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var pageHTML string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://completionator.com/Game?keyword="+query+"&sortColumn=GameName&sortDirection=ASC"),
		chromedp.WaitVisible(`.cgpager-results`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &pageHTML),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatal("No such game found within the timeout of 3 seconds!")
		} else {
			log.Fatal(err)
		}
	}

	// extract the link to the first game in the list
	return extractLinkCompletionator(pageHTML)
}

func extractLinkCompletionator(pageHTML string) (gameLink string) {
	// find the location where the first item from the search list is
	firstindex := strings.Index(pageHTML, `tr class=" even"`) + 16

	// navigate to where `a href=` is in the next part of the string (link is right after this)
	ahrefindex := strings.Index(pageHTML[firstindex:], "a href=") + 8

	// cut off all trailing info (link is before this)
	suffindex := strings.Index(pageHTML[firstindex+ahrefindex:], ">") - 1

	// take the link between the two
	return pageHTML[firstindex+ahrefindex : firstindex+ahrefindex+suffindex]
}
