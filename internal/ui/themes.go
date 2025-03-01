package ui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"gopkg.in/yaml.v3"
)

// YAML file contents must match this struct
type ColorTheme struct {
	Name            string `yaml:"name"`
	Background      string `yaml:"background"`
	Foreground      string `yaml:"foreground"`
	Primary         string `yaml:"primary"`
	ButtonColor     string `yaml:"button"`
	TextColor       string `yaml:"text"`
	DisabledColor   string `yaml:"disabled"`
	PlaceholderText string `yaml:"placeholder"`
	HoverColor      string `yaml:"hover"`
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
	case theme.ColorNameDisabled:
		return hexToColor(t.colors.DisabledColor)
	case theme.ColorNamePlaceHolder:
		return hexToColor(t.colors.PlaceholderText)
	case theme.ColorNameHover:
		return hexToColor(t.colors.HoverColor)
	default:
		return t.Theme.Color(name, variant)
	}
}

// override dflt txt size
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
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
func createLDThemes(themesDir string) {
	lightTheme := ColorTheme{
		Name:            "Light",
		Background:      "#ffffff",
		Foreground:      "#000000",
		Primary:         "#3f51b5",
		ButtonColor:     "#2196f3",
		TextColor:       "#212121",
		DisabledColor:   "#9e9e9e",
		PlaceholderText: "#757575",
		HoverColor:      "#e0e0e0",
	}

	darkTheme := ColorTheme{
		Name:            "Dark",
		Background:      "#282828",
		Foreground:      "#ebdbb2",
		Primary:         "#d79921",
		ButtonColor:     "#98971a",
		TextColor:       "#fbf1c7",
		DisabledColor:   "#7c6f64",
		PlaceholderText: "#a89984",
		HoverColor:      "#3c3836",
	}

	// add the above defined themes to themes array
	themes := []ColorTheme{lightTheme, darkTheme}

	// for each theme, make a yaml file
	for _, theme := range themes {
		data, _ := yaml.Marshal(theme)
		os.WriteFile(filepath.Join(themesDir, theme.Name+".yaml"), data, 0644)
	}
}

// loads all theme files from dir called "themes"
// INFO: this dir is in the same dir as where export files and db are stored
func loadAllthemes(themesDir string) (map[string]ColorTheme, error) {
	themes := make(map[string]ColorTheme)

	// if dir doesnt exist then create it
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
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

	// for all files in themesDir, extract the theme as ColorTheme struct and add to map
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
