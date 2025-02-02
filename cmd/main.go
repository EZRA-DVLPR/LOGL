package main

import (
	"fmt"

	"github.com/EZRA-DVLPR/GameList/internal/database"
	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

func main() {
	fmt.Println("Program Start")

	database.CreateDB()

	data := scraper.FetchHLTB("https://howlongtobeat.com/game/42069")
	database.AddToDB(data)

	database.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/155106")
	database.AddToDB(data)

	database.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/135862")
	database.AddToDB(data)

	database.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/80199")
	database.AddToDB(data)

	database.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/2127")
	database.AddToDB(data)

	database.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/68151")
	database.AddToDB(data)

	database.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/116471")
	database.AddToDB(data)

	database.PrintAllGames()

	database.DeleteFromDB(data)

	database.PrintAllGames()
}
