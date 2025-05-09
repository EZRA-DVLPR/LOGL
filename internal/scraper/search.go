package scraper

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// given a name for a game, returns the link for the game
// eg. /game/68151
func searchHLTB(query string) (gameLink string) {
	// Define a custom user agent
	userAgent := getRandUserAgent()

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
	ctx, cancel := context.WithTimeout(allocCtx, 3*time.Second)
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
			log.Println("No such game found within the timeout of 3 seconds!")
			return ""
		} else {
			log.Println(err)
			return ""
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
	userAgent := getRandUserAgent()

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
			log.Println("No such game found within the timeout of 3 seconds!")
			return ""
		} else {
			log.Println(err)
			return ""
		}
	}

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

// searches Bing for game that failed HLTB query
func searchBing(query string) (gameLink string) {
	userAgent := getRandUserAgent()

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
		// make bing search
		chromedp.Navigate("https://www.bing.com/search?q=hltb+"+query),

		// wait 0.5 seconds to grab data
		chromedp.Sleep(500*time.Millisecond),

		chromedp.OuterHTML("html", &pageHTML),
	)
	if err != nil {
		log.Println("err", err)
		return ""
	}

	return extractLinkBing(pageHTML)
}

// given the pageHTML, will scrape the first link and assume that it leads to the correct location for the correct game
func extractLinkBing(pageHTML string) (gameLink string) {
	// find the location where the first item from the search list is
	firstindex := strings.Index(pageHTML, `ol id="b_results"`) + 17

	// if there is no such item, then return empty string
	if firstindex == -1 {
		return ""
	}

	// INFO: if deeplinks does not exist, then must trim the bad stuff
	// eg. `/completions`

	// grab the link within the `href="` and `"`
	ahrefindexstart := strings.Index(pageHTML[firstindex:], `href="`) + 6
	ahrefindexend := strings.Index(pageHTML[firstindex+ahrefindexstart:], `"`)

	gameLink = pageHTML[firstindex+ahrefindexstart : firstindex+ahrefindexstart+ahrefindexend]

	// use REGEX to strip the link of what is undesired
	// `https://howlongtobeat.com/game/####` is the desired string

	pattern := `https:\/\/howlongtobeat\.com\/game\/[0-9]+`
	re := regexp.MustCompile(pattern)
	gameLink = re.FindString(gameLink)

	return gameLink
}

func getRandUserAgent() string {
	userAgents := []string{
		// chrome
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux i686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",

		// firefox
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:119.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:118.0) Gecko/20100101 Firefox/118.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3; rv:119.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6; rv:118.0) Gecko/20100101 Firefox/118.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (X11; Linux i686; rv:118.0) Gecko/20100101 Firefox/118.0",

		// safari
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Version/17.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Version/16.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Version/15.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/537.36",
		"Mozilla/5.0 (iPad; CPU OS 16_2 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_1 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/537.36",

		// edge
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",

		// mobile
		"Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; Samsung Galaxy S21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; OnePlus 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/537.36",
	}
	return userAgents[rand.Intn(len(userAgents))]
}
