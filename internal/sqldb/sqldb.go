package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
	_ "github.com/mattn/go-sqlite3"
)

// INFO: STRUCTURE OF THE DB
// games {
// 		name		PRIMARY KEY
// 		url
//		favorite
//		main
//		mainPlus
//		comp
//	}

func CreateDB() {
	fmt.Println("Creating the DB")

	// create/open the db file
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal(err)
	}

	// close the connection when done
	defer db.Close()

	// create the tables if they don't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS games (
		name TEXT PRIMARY KEY,
		url TEXT,
		favorite INTEGER,
		main TEXT,
		mainPlus TEXT,
		comp TEXT
	);
	`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	fmt.Println("Created the local DB successfully")
}

func DeleteFromDB(game scraper.Game) {
	// find the game name
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Failed to access db")
	}

	// delete from the games table
	res, err := db.Exec("DELETE FROM games WHERE name = ?", game.Name)
	if err != nil {
		log.Fatal("Error deleting game from games table: ", err)
	}

	// check if there was a deletion
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error checking affected rows: ", err)
	}
	if rowsAffected == 0 {
		fmt.Printf("Game `%s` not found in local database\n", game.Name)
		return
	}

	fmt.Println("Game deleted: ", game.Name)
}

func AddToDB(game scraper.Game) {
	// if the given game is not empty, then add to the DB
	if (game.Name == "") &&
		(game.Url == "") &&
		(game.Main == "") &&
		(game.MainPlus == "") &&
		(game.Comp == "") {
		fmt.Println("No game data received for associate game.")
		return
	}

	// open the DB
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Failed to access db")
	}
	defer db.Close()

	// if the given game already exists in the db, then dont add it
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM games WHERE name = ?)", game.Name).Scan(&exists)
	if err != nil {
		log.Fatal("Error checking game existence", err)
	}
	if exists {
		fmt.Println("Game already exists in local DB!\nSkipping insertion")
		return
	}

	fmt.Println("Adding the game data to the local DB")

	_, err = db.Exec("INSERT OR IGNORE INTO games (name, url, favorite, main, mainPlus, comp) VALUES (?,?,?,?,?,?)", game.Name, game.Url, game.Favorite, game.Main, game.MainPlus, game.Comp)
	if err != nil {
		log.Fatal("Error inserting game: ", err)
	}

	fmt.Println("Finished adding the game data to the local DB")
}

func PrintAllGames() {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	// TODO:Handle the case in which the db is empty to indicate it is empty

	rows, err := db.Query("SELECT * FROM games")
	if err != nil {
		log.Fatal("Error retrieving games: ", err)
	}
	defer rows.Close()

	fmt.Println("Games in DB:")

	for rows.Next() {
		var name, url, main, mainPlus, comp string
		var favorite int
		if err := rows.Scan(&name, &url, &favorite, &main, &mainPlus, &comp); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		fmt.Printf("Name: %s\nURL: %s\nFavorite: %d\nMain:\t%s\nMain+:\t%s\nComp:\t%s\n", name, url, favorite, main, mainPlus, comp)
	}
}
