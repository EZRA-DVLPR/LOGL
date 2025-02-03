package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
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

	// for each row in games, add a line in the markdown file
	for rows.Next() {
		var name, url string
		if err := rows.Scan(&name, &url); err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		fmt.Printf("Name: %s \nURL %s,\n", name, url)

		timesrows, err := db.Query("SELECT label, length FROM times WHERE game_name = ?", name)
		if err != nil {
			log.Fatal("Error retrieving times: ", err)
		}
		defer timesrows.Close()

		for timesrows.Next() {
			var label, length string
			if err := timesrows.Scan(&label, &length); err != nil {
				log.Fatal("Error scanning timesrow: ", err)
			}
			fmt.Printf("\t%s:\t%s\n", label, length)
		}

		// | name | length (Main Story) | length (Main + Sides) | length (Completionist)

	}
}
