package tests

import (
	"reflect"
	"testing"

	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

func TestMyFunction(t *testing.T) {
	actual := scraper.FetchHLTB("https://howlongtobeat.com/game/42069")

	expected := scraper.Game{
		Name:    "Lara Craft Go",
		Url:     "https://howlongtobeat.com/game/42069",
		Labels:  []string{"Main Story", "Main + Sides", "Completionist", "All Styles"},
		Lengths: []string{"--", "--", "--", "--"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}
