package main

import (
	"fmt"

	"github.com/EZRA-DVLPR/GameList/internal/database"
	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

func main() {
	fmt.Println("Program Start")

	// database.CreateDB()
	//
	data := scraper.FetchHLTB("https://howlongtobeat.com/game/42069")
	//
	// database.AddToDB(data)
	//
	database.PrintAllGames()
	//
	database.DeleteFromDB(data)

	database.PrintAllGames()
}
