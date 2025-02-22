package dbhandler

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
	_ "github.com/mattn/go-sqlite3"
)

// INFO: STRUCTURE OF THE DB
// games (table) {
// 		name					string		(text)			PRIMARY KEY
// 		hltburl					string		(text)
// 		completionatorurl		string		(text)
//		favorite				int			(integer)
//		main					float		(real)
//		mainPlus				float		(real)
//		comp					float		(real)
//		}

// creates the DB with table
func CreateDB() {
	fmt.Println("Creating the DB")

	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS games (
		name TEXT PRIMARY KEY,
		hltburl TEXT,
		completionatorurl TEXT,
		favorite INTEGER,
		main REAL,
		mainPlus REAL,
		comp REAL
	);
	`)
	if err != nil {
		log.Fatal("Error creating games table:", err)
	}

	fmt.Println("Created the local DBs successfully")
}

func CheckDBExists() bool {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name = ?", "games").Scan(&name)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatal("Error checking table with no data:", err)
	}
	return true
}

// given a game struct, will search DB for the name of the game to delete it
func DeleteFromDB(gamename string) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Failed to access db")
	}

	res, err := db.Exec("DELETE FROM games WHERE name = ?", gamename)
	if err != nil {
		log.Fatal("Error deleting game from games table: ", err)
	}

	if rowsAffected(res, gamename) {
		fmt.Println("Game deleted: ", gamename)
	}
}

// if the given game is not empty and not already existent in DB, then add to the DB
func AddToDB(game scraper.Game) {
	if (game.Name == "") &&
		(game.HLTBUrl == "") &&
		(game.CompletionatorUrl == "") &&
		(game.Main == -1) &&
		(game.MainPlus == -1) &&
		(game.Comp == -1) {
		fmt.Println("No game data received for associate game.")
		return
	}

	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Failed to access db")
	}
	defer db.Close()

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

	_, err = db.Exec(
		"INSERT OR IGNORE INTO games (name, hltburl, completionatorurl, favorite, main, mainPlus, comp) VALUES (?,?,?,?,?,?)",
		game.Name,
		game.HLTBUrl,
		game.CompletionatorUrl,
		game.Favorite,
		game.Main,
		game.MainPlus,
		game.Comp,
	)
	if err != nil {
		log.Fatal("Error inserting game: ", err)
	}

	fmt.Println("Finished adding the game data to the local DB")
}

// if the given game is not empty, then toggle favorite
func ToggleFavorite(gamename string) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Failed to access db")
	}
	defer db.Close()

	// get value of favorite for given game
	var favorite bool
	err = db.QueryRow("SELECT favorite FROM games WHERE name = ?", gamename).Scan(&favorite)
	if err != nil {
		log.Fatal("Error obtaining favorite value from game", err)
	}

	// update game favorite value to the opposite value
	res, err := db.Exec("UPDATE games SET favorite = ? WHERE name = ?", !favorite, gamename)
	if err != nil {
		log.Fatal("Error updating game to be favorite", err)
	}

	if rowsAffected(res, gamename) {
		fmt.Println("Toggled Favorite for given game:", gamename)
	}
}

// returns query from db as [][]string given cat, ord, and query
func SortDB(sortCategory string, sortOrder bool, queryname string) (dbOutput [][]string) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	// if sortOrder is true => ASC. false => DESC
	so := ""
	if sortOrder {
		so = "ASC"
	} else {
		so = "DESC"
	}

	// if queryname is empty, sort DB without searching for similar game names
	var rows *sql.Rows
	if queryname == "" {
		rows, err = db.Query(
			fmt.Sprintf(`
				SELECT name, main, mainPlus, comp 
				FROM games 
				ORDER BY favorite DESC, 
				CASE 
					WHEN typeof(%[1]s) = 'integer' OR %[1]s GLOB '[0-9]*' THEN CAST(%[1]s AS INTEGER) 
					ELSE %[1]s 
				END %[2]s;`,
				sortCategory,
				so,
			),
		)
	} else {
		rows, err = db.Query(
			fmt.Sprintf(`
				SELECT name, main, mainPlus, comp 
				FROM games 
				WHERE name LIKE ?
				ORDER BY favorite DESC, 
				CASE 
					WHEN typeof(%[1]s) = 'integer' OR %[1]s GLOB '[0-9]*' THEN CAST(%[1]s AS INTEGER) 
					ELSE %[1]s 
				END %[2]s;`,
				sortCategory,
				so,
			), "%"+queryname+"%")
	}
	if err != nil {
		log.Fatal("Error sorting games from games table: ", err)
	}

	// format data for return
	for rows.Next() {
		var name string
		var main, mainPlus, comp float64
		if err := rows.Scan(&name, &main, &mainPlus, &comp); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		dbOutput = append(dbOutput, []string{
			name,
			strconv.FormatFloat(main, 'f', -1, 64),
			strconv.FormatFloat(mainPlus, 'f', -1, 64),
			strconv.FormatFloat(comp, 'f', -1, 64),
		})
	}
	return dbOutput
}

func convertRowToInterface(row []string) []interface{} {
	result := make([]interface{}, len(row))
	for i, v := range row {
		result[i] = v
	}
	return result
}

func join(elements []string, sep string) string {
	if len(elements) == 0 {
		return ""
	}
	result := elements[0]
	for _, element := range elements[1:] {
		result += sep + element
	}
	return result
}

// if given rows were affected then returns true. o/w false
func rowsAffected(res sql.Result, name string) (wereAffected bool) {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error checking affected rows: ", err)
	}
	if rowsAffected == 0 {
		fmt.Printf("Game `%s` not found in local database\n", name)
		return false
	}
	return true
}
