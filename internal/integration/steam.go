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

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// timeout for 10 seconds
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
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

		// wait for js to load the games list (5 seconds)
		chromedp.Sleep(5*time.Second),

		// ensure list is present
		chromedp.WaitVisible(gamesRootSelector, chromedp.ByQuery),

		// extract games
		chromedp.Evaluate(`Array.from(document.querySelectorAll('`+gameLinksSelector+`')).map(el => el.textContent.trim())`, &gameNames))
	if err != nil {
		log.Fatal("Failed to fetch game names:", err)
	}

	fmt.Println("Owned Games:")
	for _, name := range gameNames {
		fmt.Println(name)
	}
}
