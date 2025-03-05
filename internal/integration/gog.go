package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

type GOGPage struct {
	SortBy                     any          `json:"-"`
	Page                       int          `json:"page"` // used for checking current page
	TotalProducts              any          `json:"-"`
	TotalPages                 int          `json:"totalPages"` // check bounds for number of pages
	ProductsPerPage            any          `json:"-"`
	ContentSystemCompatibility any          `json:"-"`
	Tags                       any          `json:"-"`
	Products                   []GOGProduct `json:"products"` // the products array will have the titles. i.e. the data i want to extract
	UpdatedProductsCount       any          `json:"-"`
	HiddenUpdatedProductsCount any          `json:"-"`
	AppliedFilters             any          `json:"-"`
	HasHiddenProducts          any          `json:"-"`
}

type GOGProduct struct {
	IsGalaxyCompatible   any    `json:"-"`
	Tags                 any    `json:"-"`
	ID                   any    `json:"-"`
	Availability         any    `json:"-"`
	Title                string `json:"title"` // what i want to extract
	Image                any    `json:"-"`
	URL                  any    `json:"-"`
	WorksOn              any    `json:"-"`
	Category             any    `json:"-"`
	Rating               any    `json:"-"`
	IsComingSoon         any    `json:"-"`
	IsMovie              any    `json:"-"`
	IsGame               any    `json:"-"`
	Slug                 any    `json:"-"`
	Updates              any    `json:"-"`
	IsNew                any    `json:"-"`
	DLCCount             any    `json:"-"`
	ReleaseDate          any    `json:"-"`
	IsBaseProductMissing any    `json:"-"`
	IsHidingDisabled     any    `json:"-"`
	IsInDevelopment      any    `json:"-"`
	ExtraInfo            any    `json:"-"`
	IsHidden             any    `json:"-"`
}

func GetAllGamesGOG(cookie string, searchSource string) {
	fmt.Println("Getting products from GOG...")

	url := "https://embed.gog.com/account/getFilteredProducts?mediaType=1&page=1"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// TEST: gog_us is the name of the cookie i have (US atm)
	// what about other nations/countries/regions? eg. EU? Oceania? Asia?
	req.Header.Set("Cookie", fmt.Sprintf("gog_us=%s", cookie))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Referer", "https://embed.gog.com/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest") // tells server it's an AJAX request

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	// grab the data from the page and parse using json
	var gogpage GOGPage
	err = json.Unmarshal(body, &gogpage)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	// from the 1st page, get the list of game titles
	var gameList []string
	for _, gogproduct := range gogpage.Products {
		gameList = append(gameList, gogproduct.Title)
	}

	// for each page in range [2:totalPages] inclusive, want to grab all games from each page
	for i := 2; i <= gogpage.TotalPages; i++ {
		// unpack the elts from the search from each page and append to gameList
		gameList = append(gameList, getGOGGames(i, cookie)...)
	}

	// we now have the entire list of games
	for _, game := range gameList {
		fmt.Println(game)
		dbhandler.SearchAddToDB(game, searchSource)
	}
}

func getGOGGames(pagenumber int, cookie string) (gameList []string) {
	url := fmt.Sprintf("https://embed.gog.com/account/getFilteredProducts?mediaType=1&page=%d", pagenumber)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// TEST: gog_us is the name of the cookie i have (US atm)
	// what about other nations/countries/regions? eg. EU? Oceania? Asia?
	req.Header.Set("Cookie", fmt.Sprintf("gog_us=%s", cookie))

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Referer", "https://embed.gog.com/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	// grab the data from the page and parse using json
	var gogpage GOGPage
	err = json.Unmarshal(body, &gogpage)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	// from the current page, get the list of game titles
	for _, gogproduct := range gogpage.Products {
		gameList = append(gameList, gogproduct.Title)
	}
	return
}
