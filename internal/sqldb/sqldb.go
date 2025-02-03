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
// 		name
// 		url }
// times {
// 		game_name
// 		main
//		mainPlus
// 		comp }
// where game_name is the same (references) the name of the game in the games table

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
		url TEXT
	);
	CREATE TABLE IF NOT EXISTS times (
		game_name TEXT,
		main TEXT,
		mainPlus TEXT,
		comp TEXT,
		PRIMARY KEY (game_name),
		FOREIGN KEY (game_name) REFERENCES games(name) ON DELETE CASCADE
	);
	`)
	if err != nil {
		log.Fatal("Error creating tables: ", err)
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
		fmt.Printf("Game `%s` not found in database\n", game.Name)
		return
	}

	res, err = db.Exec("DELETE FROM times WHERE game_name = ?", game.Name)
	if err != nil {
		log.Fatal("Error deleting game from times table: ", err)
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
		fmt.Println("Game already exists in database!\nSkipping insertion")
		return
	}

	fmt.Println("Adding the game data to the DB")

	// insert into games table: name, url
	_, err = db.Exec("INSERT OR IGNORE INTO games (name, url) VALUES (?,?)", game.Name, game.Url)
	if err != nil {
		log.Fatal("Error inserting game: ", err)
	}

	// insert into times table: Main, MainPlus, Comp based on game name associated
	_, err = db.Exec("INSERT INTO times (game_name, main, mainPlus, comp) VALUES (?, ?, ?, ?)", game.Name, game.Main, game.MainPlus, game.Comp)
	if err != nil {
		log.Fatal("Error inserting times", game.Name, game.Main, game.MainPlus, game.Comp)
	}

	fmt.Println("Finished adding the game data to the DB")
}

func PrintAllGames() {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	// TODO:Handle the case in which the db is empty to indicate it is empty

	rows, err := db.Query("SELECT name, url FROM games")
	if err != nil {
		log.Fatal("Error retrieving games: ", err)
	}
	defer rows.Close()

	fmt.Println("Games in DB:")

	for rows.Next() {
		var name, url string
		if err := rows.Scan(&name, &url); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		fmt.Printf("Name: %s \nURL %s,\n", name, url)

		timesrows, err := db.Query("SELECT main, mainPlus, comp FROM times WHERE game_name = ?", name)
		if err != nil {
			log.Fatal("Error retrieving times: ", err)
		}
		defer timesrows.Close()

		for timesrows.Next() {
			var main, mainPlus, comp string
			if err := timesrows.Scan(&main, &mainPlus, &comp); err != nil {
				log.Fatal("Error scanning timesrow: ", err)
			}
			fmt.Printf("\t%s:\t%s\t%s\n", main, mainPlus, comp)
		}
	}
}
