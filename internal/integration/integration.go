package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GOGPage struct {
	SortBy                     string        `json:"sortBy"`
	Page                       int           `json:"page"`
	TotalProducts              int           `json:"totalProducts"`
	TotalPages                 int           `json:"totalPages"`
	ProductsPerPage            int           `json:"productsPerPage"`
	ContentSystemCompatibility string        `json:"contentSystemCompatibility"`
	Tags                       []interface{} `json:"tags"`
	Products                   []GOGProduct  `json:"products"`
	UpdatedProductsCount       int           `json:"updatedProductsCount"`
	HiddenUpdatedProductsCount int           `json:"hiddenUpdatedProductsCount"`
	AppliedFilters             []interface{} `json:"appliedFilters"`
	HasHiddenProducts          bool          `json:"hasHiddenProducts"`
}

type GOGProduct struct {
	IsGalaxyCompatible   bool          `json:"isGalaxyCompatible"`
	Tags                 []string      `json:"tags"`
	ID                   int           `json:"id"`
	Availability         []interface{} `json:"availability"`
	Title                string        `json:"title"` // what i want to extract
	Image                string        `json:"image"`
	URL                  string        `json:"url"`
	WorksOn              []interface{} `json:"worksOn"`
	Category             string        `json:"category"`
	Rating               int           `json:"rating"`
	IsComingSoon         bool          `json:"isComingSoon"`
	IsMovie              bool          `json:"isMovie"`
	IsGame               bool          `json:"isGame"`
	Slug                 string        `json:"slug"`
	Updates              int           `json:"updates"`
	IsNew                bool          `json:"isNew"`
	DLCCount             int           `json:"dlcCount"`
	ReleaseDate          []interface{} `json:"releaseDate"`
	IsBaseProductMissing bool          `json:"isBaseProductMissing"`
	IsHidingDisabled     bool          `json:"isHidingDisabled"`
	IsInDevelopment      bool          `json:"isInDevelopment"`
	ExtraInfo            []string      `json:"extraInfo"`
	IsHidden             bool          `json:"isHidden"`
}

func GetAllGamesGOG(temp string) {
	fmt.Println("Getting products from GOG...")

	url := "https://embed.gog.com/account/getFilteredProducts?mediaType=1&page=1"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// authentication cookie
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// TEST: gog_us is the name of the cookie i have (US based btw)
	// what about other nations/countries/regions? eg. EU? Oceania? Asia?
	req.Header.Set("Cookie", fmt.Sprintf("gog_us=%s", os.Getenv("gogcookie")))

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

	fmt.Println(string(body))
	// TODO: need to get the total number of pages and repeat process for all pages
	//

	// obtain the title from the previously defined JSON structs
}
