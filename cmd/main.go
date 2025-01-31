package main

import (
	"fmt"

	// "github.com/EZRA-DVLPR/GameList/internal/scraper"
	"github.com/EZRA-DVLPR/GameList/internal/database"
)

func main() {
	fmt.Println("Program Start")

	database.CreateDB()
}
