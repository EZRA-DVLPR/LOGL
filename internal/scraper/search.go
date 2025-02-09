package scraper

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
)

func SearchHLTB(query string) (gameLink string) {
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

	// Perform the search on HLTB
	var pageHTML string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.howlongtobeat.com/?q="+query),
		chromedp.OuterHTML("html", &pageHTML),
	)
	// std logging
	if err != nil {
		log.Fatal(err)
	}

	// extract the link to the first game in the list
	return ExtractLink(pageHTML)
}

func ExtractLink(pageHTML string) (gameLink string) {
	// find the location where the first item from the search list is
	firstindex := strings.Index(pageHTML, "GameCard_inside_blur__cP8_l")

	// if there is no such item, then return
	if firstindex == -1 {
		return "No games found for given title!"
	}

	// navigate to where `a href=` is in the next part of the string
	ahrefindex := strings.Index(pageHTML[firstindex:], "a href=")

	// cut off all trailing info
	// we add 7 to cut off the `a href=` from the search
	suffindex := strings.Index(pageHTML[firstindex+ahrefindex+7:], ">")

	// grab the first link for game data that we find
	// we add 8 to cut off the `a href="` from the search
	// we also add 7 for the same reason as above, but also subtract 1
	// so that it can remove the `"`
	return pageHTML[firstindex+ahrefindex+8 : firstindex+ahrefindex+7+suffindex-1]
}
