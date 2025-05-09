package integration

import (
	"encoding/json"
	"log"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/model"
)

type EPICPage struct {
	Success any      `json:"-"`
	Data    EPICData `json:"data"` // where the data exists
}

type EPICData struct {
	Applications       []EPICGame `json:"applications"` // applications array holds the games
	LegacyApplications any        `json:"-"`
}

type EPICGame struct {
	CreatedAt       any    `json:"-"`
	ApplicationName string `json:"applicationName"` // data i want to extract
	ApplicationID   any    `json:"-"`
	PrivacyPolicy   any    `json:"-"`
	Logo            any    `json:"-"`
}

func GetAllGamesEpicString(input string) {
	log.Println("Getting products from Epic Games string")

	var epicpage EPICPage
	err := json.Unmarshal([]byte(input), &epicpage)
	if err != nil {
		log.Fatal("Error decoding json", err)
	}

	log.Println("Reading json input from Epic Games")
	model.SetMaxProcesses(len(epicpage.Data.Applications))
	for _, app := range epicpage.Data.Applications {
		log.Println("Game found:", app.ApplicationName)
		dbhandler.SearchAddToDB(app.ApplicationName)
	}
	log.Println("Finished adding game data from Epic Games")
}
