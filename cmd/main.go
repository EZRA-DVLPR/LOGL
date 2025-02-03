package main

import (
	"fmt"

	"github.com/EZRA-DVLPR/GameList/internal/mkdown"
	"github.com/EZRA-DVLPR/GameList/internal/scraper"
	"github.com/EZRA-DVLPR/GameList/internal/sqldb"
)

func main() {
	fmt.Println("Program Start")

	sqldb.CreateDB()
	sqldb.PrintAllGames()

	data := scraper.FetchHLTB("https://howlongtobeat.com/game/155106")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/135862")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/80199")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/2127")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/68151")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/116471")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	sqldb.DeleteFromDB(data)

	data = scraper.FetchHLTB("https://howlongtobeat.com/game/42069")
	sqldb.AddToDB(data)

	sqldb.PrintAllGames()

	mkdown.WriteToMarkdown()
}
