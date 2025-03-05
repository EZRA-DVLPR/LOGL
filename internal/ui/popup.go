package ui

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/internal/integration"
	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

// window for popup that will be modified for the following functions
var w2 fyne.Window

// confirmation window for updating/deleting all db data
var w3 fyne.Window

func singleGameNameSearchPopup(
	a fyne.App,
	searchSource binding.String,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) {
	// if w2 already exists then focus it and complete task
	if w2 != nil {
		w2.RequestFocus()
		return
	}

	// define w2 properties
	w2 = a.NewWindow("Single Game Name Search")
	w2.Resize(fyne.NewSize(400, 80))

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Game Name to Search")
	w2.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			entry,
			widget.NewButton("Begin Search", func() {
				// if entry is non-empty then perform search
				if strings.TrimSpace(entry.Text) != "" {
					log.Println("Search for game data beginning!")
					ss, _ := searchSource.Get()
					// search game data then add to db
					dbhandler.SearchAddToDB(entry.Text, ss)

					// update dbData
					updateDBData(sortCategory, sortOrder, searchText, dbData)
					forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
				} else {
					log.Println("No game name given")
				}
				w2.Close()
			}),
		),
	)
	w2.SetOnClosed(func() {
		w2 = nil
	})
	w2.Show()
}

func manualEntryPopup(
	a fyne.App,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) {
	if w2 != nil {
		w2.RequestFocus()
		return
	}

	w2 = a.NewWindow("Manual Game Data Entry")
	w2.Resize(fyne.NewSize(400, 100))

	gamename := widget.NewEntry()
	main := widget.NewEntry()
	mainplus := widget.NewEntry()
	comp := widget.NewEntry()
	hltbURL := widget.NewEntry()
	completionatorURL := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "Game Name",
				Widget: gamename,
			},
			{
				Text:   "Main (Hours)",
				Widget: main,
			},
			{
				Text:   "Main Plus Sides (Hours)",
				Widget: mainplus,
			},
			{
				Text:   "Completionist (Hours)",
				Widget: comp,
			},
			{
				Text:   "URL for HowLongToBeat",
				Widget: hltbURL,
			},
			{
				Text:   "URL for Completionator",
				Widget: completionatorURL,
			},
		},
		OnSubmit: func() {
			if strings.TrimSpace(gamename.Text) == "" ||
				strings.TrimSpace(main.Text) == "" ||
				strings.TrimSpace(mainplus.Text) == "" ||
				strings.TrimSpace(comp.Text) == "" {
				log.Println("Not enough game data given. Fill out top 4 fields")
			} else {
				if hltbURL.Text == "" {
					log.Println("No HLTB URL given for manual entry for game", strings.TrimSpace(gamename.Text))
				}
				if completionatorURL.Text == "" {
					log.Println("No Completionator URL given for manual entry for game", strings.TrimSpace(gamename.Text))
				}

				// check if main, mainplus, comp are valid floats
				mainfl, err := strconv.ParseFloat(main.Text, 64)
				if err != nil {
					log.Println("Improper value for Main Story. Make sure its a valid decimal.")
					w2.Close()
					return
				}
				mainplusfl, err := strconv.ParseFloat(mainplus.Text, 64)
				if err != nil {
					log.Println("Improper value for Main + Sides. Make sure its a valid decimal.")
					w2.Close()
					return
				}
				compfl, err := strconv.ParseFloat(comp.Text, 64)
				if err != nil {
					log.Println("Improper value for Completionist. Make sure its a valid decimal.")
					w2.Close()
					return
				}

				// insert the data into the db
				var newgame scraper.Game
				newgame.Name = strings.TrimSpace(gamename.Text)
				newgame.Main = float32(mainfl)
				newgame.MainPlus = float32(mainplusfl)
				newgame.Comp = float32(compfl)
				newgame.HLTBUrl = strings.TrimSpace(hltbURL.Text)
				newgame.CompletionatorUrl = strings.TrimSpace(completionatorURL.Text)
				newgame.Favorite = 0

				dbhandler.AddToDB(newgame)
				forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
			}
			w2.Close()
		},
		OnCancel: func() {
			w2.Close()
		},
	}
	w2.SetContent(
		form,
	)
	w2.SetOnClosed(func() {
		w2 = nil
	})
	w2.Show()
}

func integrationImport(
	searchSource binding.String,
	name string,
	w fyne.Window,
) {
	var main string
	var cookie string

	switch name {
	case "gog":
		main = "Enter `gog_us` cookie"
	case "psn":
		main = "Enter PSN profile name"
	case "steam":
		main = "Enter Steam profile name"
		cookie = "Enter `sessionid` cookie"
	case "epic":
		main = "Enter Epic JSON data"
	}

	// open popup with entry to receive cookie/json/profile data
	var list []*widget.FormItem
	var cookieWidget *widget.Entry

	// all options have a required main data
	mainWidget := widget.NewEntry()
	list = append(list, widget.NewFormItem(main, mainWidget))

	// steam requires the cookie as well
	if name == "steam" {
		cookieWidget = widget.NewEntry()
		list = append(list, widget.NewFormItem(cookie, cookieWidget))
	}

	dialog.ShowForm(
		"Enter Required Information Here",
		"Confirm",
		"Cancel",
		list,
		func(submitted bool) {
			if submitted {
				// check if all fields have proper information
				valid := true
				for _, ent := range list {
					if strings.TrimSpace(ent.Widget.(*widget.Entry).Text) == "" {
						valid = false
						break
					}
				}

				if valid {
					ss, _ := searchSource.Get()
					switch name {
					case "gog":
						integration.GetAllGamesGOG(mainWidget.Text, ss)
					case "psn":
						integration.GetAllGamesPS(mainWidget.Text, ss)
					case "steam":
						integration.GetAllGamesSteam(mainWidget.Text, cookieWidget.Text, ss)
					case "epic":
						integration.GetAllGamesEpicString(mainWidget.Text, ss)
					}
				} else {
					log.Println("Please ensure all fields have proper input for:", name)
				}

			} else {
				log.Println("Canceled input for:", name)
			}
		},
		w,
	)
}

func settingsPopup(
	a fyne.App,
	searchSource binding.String,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	selectedRow binding.Int,
	dbData *MyDataBinding,
	textSize binding.Float,
	selectedTheme binding.String,
) {
	if w2 != nil {
		w2.RequestFocus()
		return
	}

	w2 = a.NewWindow("Settings Window")
	w2.Resize(fyne.NewSize(400, 600))

	w2.SetContent(
		container.NewVScroll(
			container.New(
				layout.NewVBoxLayout(),
				searchSourceRadioWidget(searchSource),
				widget.NewSeparator(),
				themeSelector(selectedTheme, textSize, a),
				widget.NewSeparator(),
				textSlider(selectedTheme, textSize, a),
				widget.NewSeparator(),
				updateAllButton(a, sortCategory, sortOrder, searchText, dbData, selectedRow),
				widget.NewSeparator(),
				deleteAllButton(a, sortCategory, sortOrder, searchText, dbData, selectedRow),
			),
		),
	)
	w2.SetOnClosed(func() {
		w2 = nil
	})
	w2.Show()
}

// radio for selection of sources
func searchSourceRadioWidget(searchSource binding.String) *fyne.Container {
	label := widget.NewLabelWithStyle(
		"Search Source Selection",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	radio := widget.NewRadioGroup(
		[]string{
			"All",
			"HLTB",
			"Completionator",
		},
		func(value string) { searchSource.Set(value) },
	)

	// set default to search source saved
	ss, _ := searchSource.Get()
	radio.SetSelected(ss)

	return container.New(
		layout.NewVBoxLayout(),
		label,
		radio,
	)
}

// selector for the theme of the application
// TODO: Binding for themesDir location
func themeSelector(
	selectedTheme binding.String,
	textSize binding.Float,
	a fyne.App,
) *fyne.Container {
	st, _ := selectedTheme.Get()
	newName := abbrevName(st)
	label := widget.NewLabelWithStyle(
		fmt.Sprintf("Current Theme: %v", newName),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	// TODO: Binding for themesDir location
	availableThemes, err := loadAllThemes("themes")
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}
	themeList := container.New(
		layout.NewVBoxLayout(),
	)

	for themeName, themeColors := range availableThemes {
		newName = abbrevName(themeName)
		button := widget.NewButton(newName, func(name string, colors ColorTheme) func() {
			return func() {
				newName = abbrevName(name)
				label.SetText(fmt.Sprintf("Current Theme: %v", newName))
				selectedTheme.Set(themeName)
				ts, _ := textSize.Get()
				a.Settings().SetTheme(
					&CustomTheme{
						Theme:    theme.DefaultTheme(),
						textSize: float32(ts),
						colors:   availableThemes[themeName],
					},
				)
			}
		}(themeName, themeColors))
		themeList.Add(button)

		colorPreviews := container.New(
			layout.NewGridLayout(8),
			fixedHeightRect(hexToColor(themeColors.Background)),
			fixedHeightRect(hexToColor(themeColors.Foreground)),
			fixedHeightRect(hexToColor(themeColors.Primary)),
			fixedHeightRect(hexToColor(themeColors.ButtonColor)),
			fixedHeightRect(hexToColor(themeColors.PlaceholderText)),
			fixedHeightRect(hexToColor(themeColors.HoverColor)),
			fixedHeightRect(hexToColor(themeColors.InputBackgroundColor)),
			fixedHeightRect(hexToColor(themeColors.ScrollBarColor)),
		)
		themeList.Add(colorPreviews)
	}

	return container.New(
		layout.NewVBoxLayout(),
		label,
		themeList,
	)
}

// TODO: Binding for themesDir location
func textSlider(
	selectedTheme binding.String,
	textSize binding.Float,
	a fyne.App,
) *fyne.Container {
	// TODO: Binding for themesDir location
	availableThemes, err := loadAllThemes("themes")
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}
	label := widget.NewLabelWithStyle(
		"Change Text and Icon Size",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	ts, _ := textSize.Get()
	currSize := widget.NewLabel(fmt.Sprintf("Current size is: %v", ts))
	moveSize := widget.NewLabel("")
	moveSize.Hide()

	slider := widget.NewSliderWithData(12, 24, textSize)
	slider.OnChanged = func(res float64) {
		moveSize.Show()
		res32 := float32(res)
		moveSize.SetText(fmt.Sprintf("New Size will be: %v", res32))
	}
	slider.OnChangeEnded = func(res float64) {
		moveSize.Hide()
		res32 := float32(res)
		st, _ := selectedTheme.Get()
		currSize.SetText(fmt.Sprintf("Current size is: %v", res32))
		a.Settings().SetTheme(
			&CustomTheme{
				Theme:    theme.DefaultTheme(),
				textSize: res32,
				colors:   availableThemes[st],
			},
		)
	}
	return container.New(
		layout.NewVBoxLayout(),
		label,
		currSize,
		slider,
		moveSize,
	)
}

// PERF: use dialog confirmation
func updateAllButton(
	a fyne.App,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) *fyne.Container {
	label := widget.NewLabelWithStyle(
		"Update All Games",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	updateAll := widget.NewButton("Update All", func() {
		if w3 != nil {
			w3.RequestFocus()
			return
		}
		w3 = a.NewWindow("Confirm Deletion of entire Database")
		w3.Resize(fyne.NewSize(400, 200))
		w3.SetContent(container.New(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				w3.Close()
				w2.RequestFocus()
			}),
			widget.NewButtonWithIcon("Update all data", theme.ConfirmIcon(), func() {
				dbhandler.UpdateEntireDB()

				forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
				w3.Close()
				w2.Close()
			}),
		))
		w3.Show()
		w3.SetOnClosed(func() {
			w3 = nil
		})
		w2.SetOnClosed(func() {
			w2 = nil
		})
	})
	return container.New(
		layout.NewVBoxLayout(),
		label,
		updateAll,
	)
}

// PERF: use dialog confirmation
func deleteAllButton(
	a fyne.App,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) *fyne.Container {
	label := widget.NewLabelWithStyle(
		"Delete All Data",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	deleteAll := widget.NewButton("Delete All", func() {
		if w3 != nil {
			w3.RequestFocus()
			return
		}
		w3 = a.NewWindow("Confirm Deletion of entire Database")
		w3.Resize(fyne.NewSize(400, 200))
		w3.SetContent(container.New(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				w3.Close()
				w2.RequestFocus()
			}),
			widget.NewButtonWithIcon("DELETE ALL DATA", theme.ConfirmIcon(), func() {
				dbhandler.DeleteAllDBData()

				forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
				w3.Close()
				w2.Close()
			}),
		))
		w3.Show()
		w3.SetOnClosed(func() {
			w3 = nil
		})
		w2.SetOnClosed(func() {
			w2 = nil
		})
	})
	return container.New(
		layout.NewVBoxLayout(),
		label,
		deleteAll,
	)
}

// return a string thats abbreviated if it needs to be abbreviated
func abbrevName(name string) string {
	var returnname string
	if len(name) > 12 {
		returnname = name[:12] + "..."
	} else {
		returnname = name
	}
	return returnname
}

// create rectangles for the preview of the color theme
func fixedHeightRect(color color.Color) *canvas.Rectangle {
	rect := canvas.NewRectangle(color)
	rect.SetMinSize(fyne.NewSize(0, 40))
	return rect
}
