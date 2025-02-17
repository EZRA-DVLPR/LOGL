package dbhandler

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
	_ "github.com/mattn/go-sqlite3"
)

// INFO: STRUCTURE OF THE DB
// games (table) {
// 		name		string				PRIMARY KEY
// 		url			string
//		favorite	int
//		main		real
//		mainPlus	real
//		comp		real
//	}

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
		url TEXT,
		favorite INTEGER,
		main REAL,
		mainPlus REAL,
		comp REAL
	);
	`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	fmt.Println("Created the local DB successfully")
}

func ImportCSV() {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	file, err := os.Open("export.csv")
	if err != nil {
		log.Fatal("error opening CSV:", err)
	}
	defer file.Close()

	// read from csv and check formatting
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}
	if len(rows) < 1 {
		log.Fatal("CSV file is empty or improperly formatted")
	}

	// create the table
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", "games").Scan(&name)
	if err != nil {
		createStmt := fmt.Sprintf("CREATE TABLE games (name TEXT PRIMARY KEY, url TEXT, favorite INTEGER, main REAL, mainPlus REAL, comp REAL)")
		_, err := db.Exec(createStmt)
		if err != nil {
			log.Fatal("Error creating table:", err)
		}
		fmt.Println("Table created")
	} else {
		log.Fatal("Error with query for table creation")
	}

	// setup transaction with dummy values
	cols := rows[0]
	temp := make([]string, len(cols))
	for i := range temp {
		temp[i] = "?"
	}
	insertStmt := fmt.Sprintf("INSERT INTO games (%s) VALUES (%s);",
		join(cols, ", "), join(temp, ", "))

	// start transaction and insert data
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
	}
	for _, row := range rows[1:] {
		_, err := tx.Exec(insertStmt, convertRowToInterface(row)...)
		if err != nil {
			tx.Rollback()
			log.Fatal("Error inserting data:", err)
		}
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Import completed successfully")
}

func ImportSQL() {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlDump, err := os.ReadFile("export.sql")
	if err != nil {
		log.Fatal(err)
	}

	// perform the import (dump)
	_, err = db.Exec(string(sqlDump))
	if err != nil {
		log.Fatal("Error importing sql database:", err)
	}

	fmt.Println("SQL database imported successfully")
}

// given a game struct, will search DB for the name of the game
func DeleteFromDB(game scraper.Game) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Failed to access db")
	}

	res, err := db.Exec("DELETE FROM games WHERE name = ?", game.Name)
	if err != nil {
		log.Fatal("Error deleting game from games table: ", err)
	}

	if rowsAffected(res, game.Name) {
		fmt.Println("Game deleted: ", game.Name)
	}
}

// if the given game is not empty and not already existent in DB, then add to the DB
func AddToDB(game scraper.Game) {
	if (game.Name == "") &&
		(game.Url == "") &&
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

	_, err = db.Exec("INSERT OR IGNORE INTO games (name, url, favorite, main, mainPlus, comp) VALUES (?,?,?,?,?,?)", game.Name, game.Url, game.Favorite, game.Main, game.MainPlus, game.Comp)
	if err != nil {
		log.Fatal("Error inserting game: ", err)
	}

	fmt.Println("Finished adding the game data to the local DB")
}

// if the given game is not empty, then add favorite
func AddFavorite(game scraper.Game) {
	if (game.Name == "") &&
		(game.Url == "") &&
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

	// update given game.Favorite = 1 (true)
	res, err := db.Exec("UPDATE games SET favorite = 1 WHERE name = ?", game.Name)
	if err != nil {
		log.Fatal("Error updating game to be favorite", err)
	}

	if rowsAffected(res, game.Name) {
		fmt.Println("Favorited:", game.Name)
	}
}

// if the given game is not empty, then remove favorite
func RemoveFavorite(game scraper.Game) {
	if (game.Name == "") &&
		(game.Url == "") &&
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

	// update given game.Favorite = 0 (false)
	res, err := db.Exec("UPDATE games SET favorite = 0 WHERE name = ?", game.Name)
	if err != nil {
		log.Fatal("Error updating game to be favorite")
	}

	if rowsAffected(res, game.Name) {
		fmt.Println("Un-Favorited:", game.Name)
	}
}

// defaults to sort by name
// o/w sorts based on these criteria:
//
//	sort == name
//	sort == main
//	sort == mainPlus
//	sort == comp
//
// in all cases, it will sort the list based on favorites first, then non-favorited entries
func SortDB(sort string, sortOpt string) (dbOutput [][]string) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT name, main, mainPlus, comp FROM games ORDER BY favorite DESC, %s %s;", sort, sortOpt))
	if err != nil {
		log.Fatal("Error sorting games from games table: ", err)
	}

	// fmt.Println("Games in DB sorted by ", sort, sortOpt)

	for rows.Next() {
		var name string
		var main, mainPlus, comp float64
		if err := rows.Scan(&name, &main, &mainPlus, &comp); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		// fmt.Printf("Name: %s\nMain:\t%v\nMain+:\t%v\nComp:\t%v\n", name, main, mainPlus, comp)
		dbOutput = append(dbOutput, []string{
			name,
			strconv.FormatFloat(main, 'f', -1, 64),
			strconv.FormatFloat(mainPlus, 'f', -1, 64),
			strconv.FormatFloat(comp, 'f', -1, 64),
		})
	}
	// fmt.Println()
	return dbOutput
}

// selector for exporting
func Export(choice int) {
	switch choice {
	case 1:
		exportSQL()
	case 2:
		exportCSV()
	case 3:
		exportMarkdown()
	default:
		log.Fatal("No such export exists!")
	}
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
