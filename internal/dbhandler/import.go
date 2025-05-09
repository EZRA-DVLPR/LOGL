package dbhandler

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EZRA-DVLPR/GameList/model"
	_ "github.com/mattn/go-sqlite3"
)

// selector for importing
func Import(choice int, filename string) {
	switch choice {
	case 1:
		importCSV(filename)
	case 2:
		importSQL(filename)
	case 3:
		importTXT(filename)
	default:
		log.Fatal("No such import exists!")
	}
}

func importCSV(filename string) {
	log.Println("Importing data from CSV: ", filename)
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	file, err := os.Open(filename)
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

	model.SetMaxProcesses(len(rows))
	// create the table if it does not exist
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='games'").Scan(&exists)
	if exists != 1 {
		log.Println("Table does not exist. Creating table to insert data")
		var name string
		err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='games'").Scan(&name)
		if err != nil {
			_, err := db.Exec("CREATE TABLE games (name TEXT PRIMARY KEY, hltburl TEXT, completionatorurl TEXT, favorite INTEGER, main REAL, mainPlus REAL, comp REAL)")
			if err != nil {
				log.Fatal("Error creating table:", err)
			}
			log.Println("Table created")
		} else {
			log.Fatal("Error with query for table creation")
		}
	}

	// setup transaction with dummy values
	// INSERT OR REPLACE INTO GAMES [colname], [colname], ... VALUES ?,?,...
	cols := rows[0]
	temp := make([]string, len(cols))
	for i := range temp {
		temp[i] = "?"
	}
	insertStmt := fmt.Sprintf(
		"INSERT OR REPLACE INTO games (%s) VALUES (%s);",
		join(cols, ", "),
		join(temp, ", "),
	)

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
	}
	log.Println("Inserting data from CSV")
	// turns `INSERT OR REPLACE INTO GAMES [colname] VALUES ?`
	// into `INSERT OR REPLACE INTO GAMES name, ... VALUE gamename,...`
	// and executes transaction for each row
	for _, row := range rows[1:] {
		_, err := tx.Exec(insertStmt, convertRowToInterface(row)...)
		if err != nil {
			tx.Rollback()
			log.Fatal("Error inserting data:", err)
		}
		model.IncrementProgress()
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	log.Println("Import from CSV completed successfully")
}

func importSQL(filename string) {
	log.Println("Importing data from SQL:", filename)
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	model.SetMaxProcesses(1)

	// drop the existing tables
	log.Println("Deleting previous data")
	_, err = db.Exec("DROP TABLE IF EXISTS games;")
	if err != nil {
		log.Fatal("Error dropping tables:", err)
	}

	sqlDump, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// perform the import (dump)
	log.Println("Dumping SQL contents")
	_, err = db.Exec(string(sqlDump))
	if err != nil {
		log.Fatal("Error importing sql database:", err)
	}
	model.IncrementProgress()

	log.Println("SQL database imported successfully")
}

func importTXT(filename string) {
	log.Println("Importing data from TXT:", filename)
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("error opening txt file:", err)
	}
	defer file.Close()

	// scan file and add new obtain data and insert into gameNames []string
	var gameNames []string
	scanner := bufio.NewScanner(file)
	log.Println("Reading file for game names")
	for scanner.Scan() {
		gameNames = append(gameNames, strings.TrimSpace(scanner.Text()))
	}

	// for each game in gameNames, perform search and add to DB
	log.Println("List of game names obtained. Will now search then add each game to DB")
	model.SetMaxProcesses(len(gameNames))
	for _, game := range gameNames {
		log.Println("Obtaining Data for game", game)
		SearchAddToDB(game)
	}

	log.Println("Finished obtaining data for games in txt file")
}
