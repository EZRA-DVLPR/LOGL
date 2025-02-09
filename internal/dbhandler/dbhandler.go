package dbhandler

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

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

// TODO: refactor the code to make more simple
// 1. make a function for checking affected rows
// 2. validate that each function cannot use a helper function for simpler processing
// 3. determine if this file should be broken into several other files for simplistic overhead

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

func ImportCSV() {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	// open csv
	file, err := os.Open("export.csv")
	if err != nil {
		log.Fatal("error opening CSV:", err)
	}
	defer file.Close()

	// read from csv
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	if len(rows) < 1 {
		log.Fatal("CSV file is empty or improperly formatted")
	}

	cols := rows[0]

	// create the table
	var name string
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	err = db.QueryRow(query, "games").Scan(&name)
	if err != nil {
		createStmt := fmt.Sprintf("CREATE TABLE games (name TEXT PRIMARY KEY, url TEXT, favorite INTEGER, main TEXT, mainPlus TEXT, comp TEXT)")
		_, err := db.Exec(createStmt)
		if err != nil {
			log.Fatal("Error creating table:", err)
		}
		fmt.Println("Table created")
	} else {
		log.Fatal("Error with query for table creation")
	}

	temp := make([]string, len(cols))
	for i := range temp {
		temp[i] = "?"
	}

	insertStmt := fmt.Sprintf("INSERT INTO games (%s) VALUES (%s);",
		join(cols, ", "), join(temp, ", "))

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
	}

	// insert the data
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

	// read the "export.sql" file
	sqlDump, err := os.ReadFile("export.sql")
	if err != nil {
		log.Fatal(err)
	}

	// perform the import
	_, err = db.Exec(string(sqlDump))
	if err != nil {
		log.Fatal("Error importing sql database:", err)
	}

	fmt.Println("SQL database imported successfully")
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

func AddFavorite(game scraper.Game) {
	// if the given game is not empty, then add favorite
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

	res, err := db.Exec("UPDATE games SET favorite = 1 WHERE name = ?", game.Name)
	if err != nil {
		log.Fatal("Error updating game to be favorite", err)
	}

	// check if there was a change
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error checking affected rows: ", err)
	}
	if rowsAffected == 0 {
		fmt.Printf("Game `%s` not found in local database\n", game.Name)
		return
	}
}

func RemoveFavorite(game scraper.Game) {
	// if the given game is not empty, then add favorite
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

	res, err := db.Exec("UPDATE games SET favorite = 0 WHERE name = ?", game.Name)
	if err != nil {
		log.Fatal("Error updating game to be favorite")
	}

	// check if there was a change
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error checking affected rows: ", err)
	}
	if rowsAffected == 0 {
		fmt.Printf("Game `%s` not found in local database\n", game.Name)
		return
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
func SortDB(sort string, sortOpt string) {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	// sort by sortCat

	// TODO: handle case where DB is empty
	rows, err := db.Query(fmt.Sprintf("SELECT name, main, mainPlus, comp FROM games ORDER BY favorite DESC, %s %s;", sort, sortOpt))
	if err != nil {
		log.Fatal("Error sorting games from games table: ", err)
	}

	fmt.Println("Games in DB sorted by ", sort, sortOpt)

	for rows.Next() {
		var name, main, mainPlus, comp string
		if err := rows.Scan(&name, &main, &mainPlus, &comp); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		fmt.Printf("Name: %s\nMain:\t%s\nMain+:\t%s\nComp:\t%s\n", name, main, mainPlus, comp)
	}
	fmt.Println()
}

func Export(choice int) {
	switch choice {
	case 1:
		exportSQL()
	case 2:
		exportCSV()
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

func exportSQL() {
	fmt.Println("Exporting to SQL file")

	outputFile := "export.sql"
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening database for copying", err)
	}
	defer db.Close()

	// open file for writing sql dump
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("Error creating SQL (dump) file:", err)
	}
	defer file.Close()

	// begin dump
	fmt.Println("Exporting database to", outputFile)
	file.WriteString("PRAGMA foreign_keys=OFF;\nBEGIN TRANSACTION;\n")

	//export scheme
	rows, err := db.Query("SELECT sql FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
	if err != nil {
		log.Fatal("Error retrieving schema:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			log.Fatal("Error scanning schema row:", err)
		}
		if schema != "" {
			file.WriteString(schema + ";\n")
		}
	}

	//export data
	tables, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
	if err != nil {
		log.Fatal("Error retrieving table names:", err)
	}
	defer tables.Close()

	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			log.Fatal("Error scanning table name:", err)
		}

		// Fetch all rows from the table
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
		if err != nil {
			log.Fatalf("Error retrieving data from %s: %v", tableName, err)
		}

		// Get column names
		cols, err := rows.Columns()
		if err != nil {
			log.Fatal("Error getting columns:", err)
		}
		numCols := len(cols)

		// Prepare for value scanning
		values := make([]interface{}, numCols)
		valuePtrs := make([]interface{}, numCols)
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Iterate through rows and generate INSERT statements
		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				log.Fatal("Error scanning row:", err)
			}

			// Convert values to SQL format
			insertValues := make([]string, numCols)
			for i, val := range values {
				switch v := val.(type) {
				case nil:
					insertValues[i] = "NULL"
				case int, int64, float64:
					insertValues[i] = fmt.Sprintf("%v", v)
				case string:
					insertValues[i] = fmt.Sprintf("'%s'", fmt.Sprintf("%s", v))
				default:
					insertValues[i] = fmt.Sprintf("'%v'", v)
				}
			}

			// Write the INSERT statement
			insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n",
				tableName,
				joinColumns(cols),
				joinColumns(insertValues))
			file.WriteString(insertStmt)
		}
		rows.Close()
	}

	// end dump
	file.WriteString("COMMIT;\nPRAGMA foreign_keys=ON;\n")
	fmt.Println("Export to SQL completed successfully.")

	return
}

func exportCSV() {
	fmt.Println("Exporting to CSV")

	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening database for export", err)
	}
	defer db.Close()

	// get all data from table
	rows, err := db.Query("SELECT * FROM games")
	if err != nil {
		log.Fatal("Error retrieving data:", err)
	}
	defer rows.Close()

	// get col names
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting column names:", err)
	}

	// open csv
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatal("Error creating csv file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write col headers
	if err := writer.Write(cols); err != nil {
		log.Fatal("Error writing CSV headers")
	}

	// write rows of data
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// scan row into value ptrs
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal("Error scanning row:", err)
		}

		// convert values to string
		stringVals := make([]string, len(cols))
		for i, val := range values {
			if val == nil {
				stringVals[i] = ""
			} else {
				stringVals[i] = fmt.Sprintf("%v", val)
			}
		}

		// write row to csv
		if err := writer.Write(stringVals); err != nil {
			log.Fatal("Error writing row to CSV:", err)
		}
	}

	fmt.Println("Export to CSV completed successfully")
}

func joinColumns(cols []string) string {
	return fmt.Sprintf("%s", join(cols, ", "))
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
