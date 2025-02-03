package mkdown

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func WriteToMarkdown() {
	db, err := sql.Open("sqlite3", "games.db")
	if err != nil {
		log.Fatal("Error accessing local dB: ", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, url FROM games")
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

	_, err = mdfile.WriteString("| **Game Name** | **Main Story** | **Main + Sides** | **Completionist** |\n")
	_, err = mdfile.WriteString("| :---- | ---- | ---- | ---- |\n")
	if err != nil {
		log.Fatal("Failed to begin writing to markdown file")
	}

	// for each row in games, add a line in the markdown file
	for rows.Next() {
		var name, url string
		if err := rows.Scan(&name, &url); err != nil {
			log.Fatal("Error scanning row: ", err)
		}

		var main, mainPlus, comp string
		err := db.QueryRow("SELECT main, mainPlus, comp FROM times WHERE game_name = ?", name).Scan(&main, &mainPlus, &comp)
		if err != nil {
			log.Fatal("Error retrieving times: ", err)
		}

		// | name | Main Story | Main + Sides | Completionist
		fmt.Println(name, main, mainPlus, comp)

		_, err = mdfile.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", name, main, mainPlus, comp))

	}
}
