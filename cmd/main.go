package main

import (
	"fmt"

	// "github.com/EZRA-DVLPR/GameList/internal/mkdown"
	// "github.com/EZRA-DVLPR/GameList/internal/scraper"
	"github.com/EZRA-DVLPR/GameList/internal/sqldb"
	// "github.com/EZRA-DVLPR/GameList/internal/ui"
)

func main() {
	fmt.Println("Program Start")

	// sqldb.CreateDB()

	//
	// data := scraper.FetchHLTB("https://howlongtobeat.com/game/155106")
	// sqldb.AddToDB(data)
	//

	//
	// data = scraper.FetchHLTB("https://howlongtobeat.com/game/135862")
	// sqldb.AddToDB(data)
	//

	//
	// data = scraper.FetchHLTB("https://howlongtobeat.com/game/80199")
	// sqldb.AddToDB(data)
	//

	//
	// data = scraper.FetchHLTB("https://howlongtobeat.com/game/2127")
	// sqldb.AddToDB(data)
	//

	//
	// data = scraper.FetchHLTB("https://howlongtobeat.com/game/68151")
	// sqldb.AddToDB(data)
	//

	//
	// data = scraper.FetchHLTB("https://howlongtobeat.com/game/116471")
	// sqldb.AddToDB(data)
	//

	//
	// sqldb.DeleteFromDB(data)
	//

	//
	// data = scraper.FetchHLTB("https://howlongtobeat.com/game/42069")
	// sqldb.AddToDB(data)
	//

	//
	// mkdown.WriteToMarkdown()
	//
	// ui.StartGUI()

	// sqldb.Export(1)
	// sqldb.Export(2)

	// sqldb.ImportSQL()
	// sqldb.ImportCSV()

	sqldb.SortDB("name", "ASC")
	sqldb.SortDB("name", "DESC")
	sqldb.SortDB("main", "ASC")
	sqldb.SortDB("main", "DESC")
}
