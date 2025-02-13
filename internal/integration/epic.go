package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

// given the data (either in string or txt file format), formats (unmarshals) into JSON which is then parsed for the desired info

func GetAllEpicGamesString(input string) {
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

func GetAllEpicGamesFile() {
	fmt.Println("Getting products from Epic Games txt file...")

	file, err := os.Open("epicgames.txt")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file", err)
	}

	var epicpage EPICPage
	err = json.Unmarshal(content, &epicpage)
	if err != nil {
		log.Fatal("Error decoding json", err)
	}

	for _, app := range epicpage.Data.Applications {
		fmt.Println(app.ApplicationName)
	}
}
