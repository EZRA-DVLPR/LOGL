package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

type Salary struct {
	Basic, HRA, TA float64
}

type Employee struct {
	FirstName, LastName, Email string
	Age                        int
	MonthlySalary              []Salary
}

func CreateDB() {
	fmt.Println("Creating the DB")

	// obtain the data for adding to the DB
	data := scraper.FetchHLTB("https://howlongtobeat.com/game/42069")

	// format the data for the file
	file, _ := json.MarshalIndent(data, "", " ")

	// write the data to the file
	_ = os.WriteFile("test.json", file, 0644)
}

func ReadDB() {
	content, err := os.ReadFile("test.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload scraper.Game
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	log.Printf("Name %s\n", payload.Name)
	log.Printf("Url %s\n", payload.Url)
	for key, val := range payload.TimeData {
		log.Println(key, val)
	}
}
