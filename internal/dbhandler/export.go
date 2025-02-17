package dbhandler

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

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

	//export schema
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

		// fetch all rows from the table
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
		if err != nil {
			log.Fatalf("Error retrieving data from %s: %v", tableName, err)
		}

		// get column names
		cols, err := rows.Columns()
		if err != nil {
			log.Fatal("Error getting columns:", err)
		}
		numCols := len(cols)

		// prepare for value scanning
		values := make([]interface{}, numCols)
		valuePtrs := make([]interface{}, numCols)
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// iterate through rows and generate INSERT statements for all values in each row
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

	// open csv to write to
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

func exportMarkdown() {
	fmt.Println("Exporting to Markdown")
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	// select everything except the url to be grabbed
	rows, err := db.Query("SELECT name, favorite, main, mainPlus, comp FROM games")
	if err != nil {
		log.Fatal("Error retrieving games: ", err)
	}
	defer rows.Close()

	// open the markdown file we are going to be writing to
	mdfile, err := os.Create("GameList.md")
	if err != nil {
		log.Fatal("Error creating markdown file", err)
	}
	defer mdfile.Close()

	_, err = mdfile.WriteString("| No. | **Game Name** | **Main Story** | **Main + Sides** | **Completionist** | Favorite |\n")
	_, err = mdfile.WriteString("| :----: | :---- | ---- | ---- | ---- | ---- |\n")
	if err != nil {
		log.Fatal("Failed to begin writing to markdown file")
	}

	id := 1
	// for each row in games, add a line in the markdown file
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

	fmt.Println("Export to Markdown completed successfully")
}

// joins columns into a single string
func joinColumns(cols []string) string {
	return fmt.Sprintf("%s", join(cols, ", "))
}
