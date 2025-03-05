package integration

import (
	"encoding/json"
	"fmt"
	"log"
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
	fmt.Println("Getting products from Epic Games string...")

	var epicpage EPICPage
	err := json.Unmarshal([]byte(input), &epicpage)
	if err != nil {
		log.Fatal("Error decoding json", err)
	}

	for _, app := range epicpage.Data.Applications {
		fmt.Println(app.ApplicationName)
	}
}
