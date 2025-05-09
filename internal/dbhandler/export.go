package dbhandler

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// selector for exporting
func Export(choice int, filename string) {
	// check filename for extension and remove it
	hasExt := strings.Index(filename, ".")
	if hasExt != -1 {
		filename = filename[:hasExt]
	}

	switch choice {
	case 1:
		exportCSV(filename)
	case 2:
		exportSQL(filename)
	case 3:
		exportMarkdown(filename)
	default:
		log.Fatal("No such export exists!")
	}
}

func exportCSV(filename string) {
	log.Println("Exporting to CSV")

	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening database for export", err)
	}
	defer db.Close()

	// get all data from table
	log.Println("Getting all game data")
	rows, err := db.Query("SELECT * FROM games")
	if err != nil {
		log.Fatal("Error retrieving data:", err)
	}
	defer rows.Close()

	// get col names
	log.Println("Getting column names")
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting column names:", err)
	}

	// open csv to write to
	file, err := os.Create(filename + ".csv")
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
	log.Println("Writing game data to csv file")
	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))
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
	log.Println("Export to CSV completed successfully")
}

func exportSQL(filename string) {
	log.Println("Exporting to SQL file")

	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error opening database for copying", err)
	}
	defer db.Close()

	// open file for writing sql dump
	file, err := os.Create(filename + ".sql")
	if err != nil {
		log.Fatal("Error creating SQL (dump) file:", err)
	}
	defer file.Close()

	// begin dump
	file.WriteString("BEGIN TRANSACTION;\n")

	//export schema
	log.Println("Extracting Schema")
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

	//export tables
	log.Println("Extracting table")
	tables, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
	if err != nil {
		log.Fatal("Error retrieving table names:", err)
	}
	defer tables.Close()

	// for each table extract all data
	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			log.Fatal("Error scanning table name:", err)
		}

		// fetch all rows from the table
		log.Println("Obtaining Row data")
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
		if err != nil {
			log.Fatalf("Error retrieving data from %s: %v", tableName, err)
		}

		// get column names
		cols, err := rows.Columns()
		log.Println("Obtaining Column data")
		if err != nil {
			log.Fatal("Error getting columns:", err)
		}
		numCols := len(cols)

		// prepare for value scanning
		values := make([]any, numCols)
		valuePtrs := make([]any, numCols)
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// iterate through rows and generate INSERT statements for all values in each row
		log.Println("Writing table contents")
		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				log.Fatal("Error scanning row:", err)
			}

			// convert values to SQL format
			insertValues := make([]string, numCols)
			for i, val := range values {
				switch v := val.(type) {
				case nil:
					insertValues[i] = "NULL"
				case int, float32:
					insertValues[i] = fmt.Sprintf("%v", v)
				case string:
					insertValues[i] = fmt.Sprintf("'%s'", fmt.Sprintf("%s", v))
				default:
					insertValues[i] = fmt.Sprintf("'%v'", v)
				}
			}

			// write the INSERT statement
			file.WriteString(
				fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) VALUES (%s);\n",
					tableName,
					joinColumns(cols),
					joinColumns(insertValues)),
			)
		}
		rows.Close()
	}

	// end dump
	file.WriteString("COMMIT;\n")
	log.Println("Export to SQL completed successfully.")

	return
}

// PERF: Export the current view, not the default one in the database
func exportMarkdown(filename string) {
	log.Println("Exporting to Markdown")
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	// select everything except the url to be grabbed
	log.Println("Obtaining Game Data")
	rows, err := db.Query("SELECT name, favorite, main, mainPlus, comp FROM games")
	if err != nil {
		log.Fatal("Error retrieving games: ", err)
	}
	defer rows.Close()

	// open the markdown file we are going to be writing to
	mdfile, err := os.Create(filename + ".md")
	if err != nil {
		log.Fatal("Error creating markdown file", err)
	}
	defer mdfile.Close()

	log.Println("Writing Headers")
	_, err = mdfile.WriteString("| No. | **Game Name** | **Main Story** | **Main + Sides** | **Completionist** | Favorite |\n")
	_, err = mdfile.WriteString("| :----: | :---- | ---- | ---- | ---- | ---- |\n")
	if err != nil {
		log.Fatal("Failed to begin writing to markdown file")
	}

	id := 1
	// for each row in games, add a line in the markdown file
	log.Println("Writing Game Data")
	for rows.Next() {
		var name string
		var main, mainPlus, comp float32
		var favorite int
		if err := rows.Scan(&name, &favorite, &main, &mainPlus, &comp); err != nil {
			log.Fatal("Error scanning row: ", err)
		}

		//  | No. | name | Main Story | Main + Sides | Completionist | Favorite |
		_, err = mdfile.WriteString(fmt.Sprintf("| %d. | %s | %v | %v | %v | %d |\n", id, name, main, mainPlus, comp, favorite))
		id += 1
	}

	log.Println("Export to Markdown completed successfully")
}

// joins columns into a single string
func joinColumns(cols []string) string {
	return fmt.Sprintf("%s", join(cols, ", "))
}
