package ui

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"gopkg.in/yaml.v3"
)

// YAML file contents must match this struct
//
// INFO: This struct affects the following fields of the theme which govern the look of the entire app:
// background: background color of the window
// altBackground: used for dbrender to show difference between odd/even rows
// foreground: text colors used for all text and buttons
// primary: settings -> highlighted (selected) src
// buttoncolor: color for buttons
// placeholderText: color of placeholder text in entry widget
// hovercolor: color when mouse cursor hovers an interactive widget
// inputbackgroundcolor: color for input fields background
// scrollbarcolor: color of scroll bar
type ColorTheme struct {
	Name                 string `yaml:"name"`
	Background           string `yaml:"background"`
	AltBackground        string `yaml:"altBackground"`
	Foreground           string `yaml:"foreground"`
	Primary              string `yaml:"primary"`
	ButtonColor          string `yaml:"button"`
	PlaceholderText      string `yaml:"placeholder"`
	HoverColor           string `yaml:"hover"`
	InputBackgroundColor string `yaml:"inputBackground"`
	ScrollBarColor       string `yaml:"scrollBar"`
}

// overwrite the theme from fyne to allow for alternate text sizes and colors
type CustomTheme struct {
	fyne.Theme
	textSize float32
	colors   ColorTheme
}

// override dflt theme colors
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return hexToColor(t.colors.Background)
	case theme.ColorNameForeground:
		return hexToColor(t.colors.Foreground)
	case theme.ColorNamePrimary:
		return hexToColor(t.colors.Primary)
	case theme.ColorNameButton:
		return hexToColor(t.colors.ButtonColor)
	case theme.ColorNamePlaceHolder:
		return hexToColor(t.colors.PlaceholderText)
	case theme.ColorNameHover, theme.ColorNameFocus:
		return hexToColor(t.colors.HoverColor)
	case theme.ColorNameInputBackground:
		return hexToColor(t.colors.InputBackgroundColor)
	case theme.ColorNameScrollBar:
		return hexToColor(t.colors.ScrollBarColor)
	default:
		return t.Theme.Color(name, variant)
	}
}

// override dflt txt size
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText || name == theme.SizeNameInlineIcon {
		return t.textSize
	}
	return t.Theme.Size(name)
}

// convert string of hex color to color.Color
// INFO: colors are in RGB format
// eg. #RRGGBB -- #010101
// which gets turned into color as RRGGBBA where A = 255
// PERF: Allow diff colors and color styles to be used
// eg. CMYK, RGBA, etc.
func hexToColor(hex string) color.Color {
	var r, g, b, a uint8 = 0, 0, 0, 255
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// loads theme from YAML file
func loadTheme(filename string) (ColorTheme, error) {
	log.Println("Loading theme from yaml file:", filename)
	var theme ColorTheme

	// if error reading file, return empty colortheme + error
	filedata, err := os.ReadFile(filename)
	if err != nil {
		return theme, err
	}

	// extract colortheme and return it + err
	err = yaml.Unmarshal(filedata, &theme)
	return theme, err
}

// makes default L/D themes as yaml files
// TODO: Make sure the themes look good and are properly described
func createLDThemes(themesDir string) {
	log.Println("Creating Light/Dark themes")
	lightTheme := ColorTheme{
		Name:                 "Light",
		Background:           "#ffffff",
		AltBackground:        "#ffffff",
		Foreground:           "#9e9e9e",
		Primary:              "#3f51b5",
		ButtonColor:          "#2196f3",
		PlaceholderText:      "#700575",
		HoverColor:           "#e0e000",
		InputBackgroundColor: "#3c3836",
		ScrollBarColor:       "#665c54",
	}

	darkTheme := ColorTheme{
		Name:                 "Dark",
		Background:           "#282828",
		AltBackground:        "#282828",
		Foreground:           "#ebdbb2",
		Primary:              "#d79921",
		ButtonColor:          "#98971a",
		PlaceholderText:      "#a89984",
		HoverColor:           "#3c3836",
		InputBackgroundColor: "#3c3836",
		ScrollBarColor:       "#665c54",
	}

	// add the above defined themes to themes array
	themes := []ColorTheme{lightTheme, darkTheme}

	log.Println("Making yaml files")
	// for each theme, make a yaml file
	for _, theme := range themes {
		data, _ := yaml.Marshal(theme)
		os.WriteFile(filepath.Join(themesDir, theme.Name+".yaml"), data, 0644)
	}
}

// loads all theme files from dir called "themes"
// INFO: this dir is in the same dir as where export files and db are stored
func loadAllThemes(themesDir string) (map[string]ColorTheme, error) {
	themes := make(map[string]ColorTheme)

	log.Println("Checking existence of given directory:", themesDir)
	// if dir doesnt exist then create it
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		log.Println("Given directory DNE. Creating it")
		err := os.MkdirAll(themesDir, 0755)
		if err != nil {
			// problem with creation so return empty map + err
			return themes, err
		}
		// create the standard L/D themes (as yaml files)
		createLDThemes(themesDir)
	}

	files, err := os.ReadDir(themesDir)
	if err != nil {
		// problem reading directory, so return empty map
		return themes, err
	}

	// check if the light/dark.yaml files exist and if not then create them
	log.Println("Locating Default themes from directory:", themesDir)
	_, errLight := os.Stat(themesDir + "/Light.yaml")
	_, errDark := os.Stat(themesDir + "/Dark.yaml")
	if (errors.Is(errLight, os.ErrNotExist)) || (errors.Is(errDark, os.ErrNotExist)) {
		log.Println("Error finding Default themes: Light and Dark. Creating them inside of :", themesDir)
		createLDThemes(themesDir)
	} else {
		log.Println("Default themes found")
	}

	log.Println("Extracting themes from themes dir")
	// for all files in themesDir, extract the theme as ColorTheme struct and add to map[string]ColorTheme
	for _, file := range files {
		if (!file.IsDir()) && (filepath.Ext(file.Name()) == ".yaml") {
			theme, err := loadTheme(filepath.Join(themesDir, file.Name()))
			if err == nil {
				themes[theme.Name] = theme
			}
		}
	}

	return themes, nil
}
