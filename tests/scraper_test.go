package tests

import (
	"testing"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

func TestScraperFetchHLTB(t *testing.T) {
	// grab a sample game with some data
	actual := scraper.FetchHLTB("https://howlongtobeat.com/game/42069")

	// explicitly write what the game data should be
	expected := scraper.Game{
		Name: "Lara Craft Go",
		Url:  "https://howlongtobeat.com/game/42069",
		TimeData: map[string]string{
			"Completionist": "--",
			"Main + Sides":  "--",
			"Main Story":    "--",
		},
	}

	// compare the values between the retrieved and hard-coded data
	if actual.Name != expected.Name {
		t.Errorf("Expected %q, got %q", expected.Name, actual.Name)
	}
	if actual.Url != expected.Url {
		t.Errorf("Expected %q, got %q", expected.Url, actual.Url)
	}
	if actual.TimeData["Completionist"] != expected.TimeData["Completionist"] {
		t.Errorf("Expected %q, got %q", expected.TimeData["Completionist"], actual.TimeData["Completionist"])
	}
	if actual.TimeData["Main + Sides"] != expected.TimeData["Main + Sides"] {
		t.Errorf("Expected %q, got %q", expected.TimeData["Main + Sides"], actual.TimeData["Main + Sides"])
	}
	if actual.TimeData["Main Story"] != expected.TimeData["Main Story"] {
		t.Errorf("Expected %q, got %q", expected.TimeData["Main Story"], actual.TimeData["Main Story"])
	}
}
