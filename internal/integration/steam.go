package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

func GetAllGamesSteam(profile string) {
	fmt.Println("Getting products from Steam...")

	url := "https://steamcommunity.com/id/" + profile + "/games/?tab-all=&tab=all"

	// load the env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// define the cookie
	cookieName := "steamLoginSecure"
	cookieValue := os.Getenv("steamloginsecure")
	cookieDomain := "steamcommunity.com"

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// copying the js code basically
	gamesRootSelector := `div[data-featuretarget="gameslist-root"]`
	gameLinksSelector := `div[data-featuretarget="gameslist-root"] div.Panel div.Panel span > a`

	var gameNames []string

	err = chromedp.Run(ctx,
		// enable network to set cookies
		network.Enable(),

		// set the cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetCookie(cookieName, cookieValue).
				WithDomain(cookieDomain).
				WithPath("/").
				Do(ctx)
		}),

		chromedp.Navigate(url),

		// Wait for JavaScript to load the games list
		chromedp.Sleep(5*time.Second),

		// Ensure the game list is present
		chromedp.WaitVisible(gamesRootSelector, chromedp.ByQuery),

		// Extract games
		chromedp.Evaluate(`Array.from(document.querySelectorAll('`+gameLinksSelector+`')).map(el => el.textContent.trim())`, &gameNames))
	if err != nil {
		log.Fatal("Failed to fetch game names:", err)
	}

	// Print the extracted game names
	fmt.Println("Owned Games:")
	for _, name := range gameNames {
		fmt.Println(name)
	}
}
