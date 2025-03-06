package integration

import (
	"context"
	"log"
	"time"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func GetAllGamesSteam(profile string, cookie string, searchSource string) {
	log.Println("Getting products from Steam for given profile:", profile)

	// define the cookie
	log.Println("Setting up HTTP request")
	cookieName := "steamLoginSecure"
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
	log.Println("Making HTTP request")
	err := chromedp.Run(ctx,
		// enable network to set cookies
		network.Enable(),

		// set the cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetCookie(cookieName, cookie).
				WithDomain(cookieDomain).
				WithPath("/").
				Do(ctx)
		}),

		chromedp.Navigate("https://steamcommunity.com/id/"+profile+"/games/?tab-all=&tab=all"),

		// wait for js to load the games list (5 seconds)
		chromedp.Sleep(5*time.Second),

		// ensure list is present
		chromedp.WaitVisible(gamesRootSelector, chromedp.ByQuery),

		// extract games
		chromedp.Evaluate(`Array.from(document.querySelectorAll('`+gameLinksSelector+`')).map(el => el.textContent.trim())`, &gameNames))
	if err != nil {
		log.Fatal("Failed to fetch game names:", err)
	}
	log.Println("HTTP Request processed successfully. List of games obtained")

	for _, name := range gameNames {
		log.Println("Game found:", name)
		dbhandler.SearchAddToDB(name, searchSource)
	}
	log.Println("Finished adding game data from Steam")
}
