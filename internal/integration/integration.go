package integration

import (
	"fmt"
	"io"
	"net/http"
)

func GetAllGamesGOG(gogcookie string) {
	fmt.Println("getting products from gog")
	url := "https://embed.gog.com/account/getFilteredProducts?mediaType=1&page=1"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	// TODO:what does the cookie look like for diff countries/continents?
	// eg. is it eu for europe?
	// gog_us is for US
	req.Header.Set("Cookie", fmt.Sprintf("gog_us=", gogcookie))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// TODO: parse the body of the response to obtain all strings with title
	// tags: -> [0-99]: -> title: -> "whatever is inside the quotes"

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
