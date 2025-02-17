package mkdown

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: bring this code into the dbhandler pkg, and make a new file just for exporting
func WriteToMarkdown() {
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
