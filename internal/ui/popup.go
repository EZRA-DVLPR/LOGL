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
	"github.com/EZRA-DVLPR/GameList/model"
)

// window for popup settings menu
var w2 fyne.Window

func singleGameNameSearchPopup(
	searchSource binding.String,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
	w fyne.Window,
) {
	var list []*widget.FormItem

	// widget to enter game name for searching
	mainWidget := widget.NewEntry()
	list = append(list, widget.NewFormItem("Game Name to Search", mainWidget))

	dialog.ShowForm(
		"Enter Game Name for Searching Here",
		"Search",
		"Cancel",
		list,
		func(submitted bool) {
			if submitted {
				// check if the game name is nonempty
				valid := true
				for _, ent := range list {
					if strings.TrimSpace(ent.Widget.(*widget.Entry).Text) == "" {
						valid = false
						break
					}
				}

				if valid {
					// bring up progress menu
					model.SetMaxProcesses(1)
					PopProgressBar(w)

					ss, _ := searchSource.Get()
					// search game data then add to db
					dbhandler.SearchAddToDB(mainWidget.Text, ss)

					// update dbData
					updateDBData(sortCategory, sortOrder, searchText, dbData)
					forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)

				} else {
					log.Println("No Game Name given for search")
				}
			} else {
				log.Println("User Cancelled Search by Game Name")
			}
		},
		w,
	)
}

func manualEntryPopup(
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
	w fyne.Window,
) {
	var list []*widget.FormItem

	// entry widgets to obtain the data from user
	gamename := widget.NewEntry()
	main := widget.NewEntry()
	mainplus := widget.NewEntry()
	comp := widget.NewEntry()
	hltbURL := widget.NewEntry()
	completionatorURL := widget.NewEntry()

	list = append(list, widget.NewFormItem("Game Name", gamename))
	list = append(list, widget.NewFormItem("Main (Hours)", main))
	list = append(list, widget.NewFormItem("Main Plus Sides (Hours)", mainplus))
	list = append(list, widget.NewFormItem("Completionist (Hours)", comp))
	list = append(list, widget.NewFormItem("URL for HowLongToBeat", hltbURL))
	list = append(list, widget.NewFormItem("URL for Completionator", completionatorURL))

	dialog.ShowForm(
		"Manually Enter the Game Data Here",
		"Manual Add",
		"Cancel",
		list,
		func(submitted bool) {
			if submitted {
				// check if entries for list are non-empty
				valid := true
				if strings.TrimSpace(gamename.Text) == "" ||
					strings.TrimSpace(main.Text) == "" ||
					strings.TrimSpace(mainplus.Text) == "" ||
					strings.TrimSpace(comp.Text) == "" {
					log.Println("Not enough game data given for manual entry. Fill out top 4 fields")
					valid = false
				}

				if valid {
					// check if main, mainplus, comp are valid floats
					mainfl, err := strconv.ParseFloat(main.Text, 64)
					if err != nil {
						log.Println("Improper value for Main Story. Make sure its a valid decimal")
						return
					}
					mainplusfl, err := strconv.ParseFloat(mainplus.Text, 64)
					if err != nil {
						log.Println("Improper value for Main + Sides. Make sure its a valid decimal")
						return
					}
					compfl, err := strconv.ParseFloat(comp.Text, 64)
					if err != nil {
						log.Println("Improper value for Completionist. Make sure its a valid decimal")
						return
					}

					// check if the URLs are given
					// NOTE: this doesn't affect inputting into DB
					if hltbURL.Text == "" {
						log.Println("No HLTB URL given for manual entry for game", strings.TrimSpace(gamename.Text))
					}
					if completionatorURL.Text == "" {
						log.Println("No Completionator URL given for manual entry for game", strings.TrimSpace(gamename.Text))
					}

					// make the Game struct then add it to the db
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
				} else {
					log.Println("No Game Name given for search")
				}
			} else {
				log.Println("User Cancelled Search by Game Name")
			}
		},
		w,
	)
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
					log.Println("Sending all fields to integration:", name)
					switch name {
					case "gog":
						integration.GetAllGamesGOG(mainWidget.Text, ss)
					case "psn":
						integration.GetAllGamesPS(mainWidget.Text, ss)
					case "steam":
						integration.GetAllGamesSteam(mainWidget.Text, cookieWidget.Text, ss)
					case "epic":
						integration.GetAllGamesEpicString(mainWidget.Text, ss)
					default:
						log.Println("Integration not found:", name)
					}
				} else {
					log.Println("Please ensure all fields have proper integration input for:", name)
				}

			} else {
				log.Println("Canceled integration import for:", name)
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
	availableThemes map[string]ColorTheme,
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
				themeSelector(sortCategory, sortOrder, searchText, selectedRow, dbData, selectedTheme, textSize, availableThemes, a),
				widget.NewSeparator(),
				textSlider(selectedTheme, textSize, availableThemes, a),
				widget.NewSeparator(),
				updateAllButton(sortCategory, sortOrder, searchText, dbData, selectedRow),
				widget.NewSeparator(),
				deleteAllButton(sortCategory, sortOrder, searchText, dbData, selectedRow),
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
		func(value string) {
			searchSource.Set(value)
			log.Println("Search Source changed to:", value)
		},
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
func themeSelector(
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	selectedRow binding.Int,
	dbData *MyDataBinding,
	selectedTheme binding.String,
	textSize binding.Float,
	availableThemes map[string]ColorTheme,
	a fyne.App,
) *fyne.Container {
	st, _ := selectedTheme.Get()
	newName := abbrevName(st)
	label := widget.NewLabelWithStyle(
		fmt.Sprintf("Current Theme: %s", newName),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	themeList := container.New(
		layout.NewVBoxLayout(),
	)

	for themeName, themeColors := range availableThemes {
		newName = abbrevName(themeName)
		button := widget.NewButton(newName, func(name string, colors ColorTheme) func() {
			return func() {
				newName = abbrevName(name)
				label.SetText(fmt.Sprintf("Current Theme: %s", newName))
				selectedTheme.Set(themeName)
				log.Println("Changed Theme to:", themeName)
				ts, _ := textSize.Get()
				a.Settings().SetTheme(
					&CustomTheme{
						Theme:    theme.DefaultTheme(),
						textSize: float32(ts),
						colors:   availableThemes[themeName],
					},
				)
				forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
			}
		}(themeName, themeColors))
		themeList.Add(button)

		colorPreviews := container.New(
			layout.NewGridLayout(9),
			fixedHeightRect(hexToColor(themeColors.Background)),
			fixedHeightRect(hexToColor(themeColors.AltBackground)),
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

func textSlider(
	selectedTheme binding.String,
	textSize binding.Float,
	availableThemes map[string]ColorTheme,
	a fyne.App,
) *fyne.Container {
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
		moveSize.SetText(fmt.Sprintf("New Size will be: %v", float32(res)))
	}
	slider.OnChangeEnded = func(res float64) {
		moveSize.Hide()
		res32 := float32(res)
		currSize.SetText(fmt.Sprintf("Current size is: %v", res32))
		log.Println(fmt.Sprintf("Text Size changed to: %v", res32))
		st, _ := selectedTheme.Get()
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

func updateAllButton(
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
		dialog.ShowConfirm(
			"Update All Game information",
			"Update All",
			func(submitted bool) {
				if submitted {
					log.Println("Updating Entire DB")
					dbhandler.UpdateEntireDB()
					forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
				}
			},
			w2,
		)
	})
	return container.New(
		layout.NewVBoxLayout(),
		label,
		updateAll,
	)
}

func deleteAllButton(
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
		dialog.ShowConfirm(
			"Delete All Game information",
			"Delete All",
			func(submitted bool) {
				if submitted {
					log.Println("Deleting data in DB")
					dbhandler.DeleteAllDBData()
					forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
				}
			},
			w2,
		)
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

func PopProgressBar(
	w fyne.Window,
) {
	// progress bar to be displayed
	progBar := widget.NewProgressBarWithData(model.GlobalModel.Progress)
	model.ResetProgress()

	//  TODO: Set max value based on the total number of things to do
	// if importing a list of games, then it should be the total number of games to import
	// if updating a list of games, then it should be the total number of games to update
	// 	if (PARAMETER TO THIS FUNCTION) {
	// 		progBar.Max = total number of games in db
	//  } else {
	// 		progBar.Max = 1
	//  }
	procmax, err := model.GetMaxProcesses()
	if err != nil {
		log.Fatal("Error getting max processes for display", err)
	}
	progBar.Max = float64(procmax)

	// create generic variable of the dialog
	var customDialog dialog.Dialog

	// create confirmation button that will close the dialog once the
	// progress bar is done. Disable the button
	actionButton := widget.NewButton("Please Wait...", func() {
		customDialog.Hide()
	})
	actionButton.Disable()

	// listener that waits for updates and increments progress
	var listener binding.DataListener

	// function for the listener to use
	listenerFunc := func(val float64) {
		// allow closing dialog and remove listener if new value is max # processes
		if val == float64(procmax) {
			actionButton.Enable()
			actionButton.SetText("Processing Completed!")
			actionButton.Refresh()
			model.RemoveProgressListener(listener)
		}
	}

	// attach progress listener
	model.AddProgressListener(listenerFunc)

	// container with the things to be displayed in the dialog
	content := container.NewVBox(
		// TODO: Change the text in the text widget based on what thing is happening
		// 	if (PARAMETER TO THIS FUNCTION) {
		// 		widget.NewLabel("Getting the new games..."),
		//  } else {
		// 		widget.NewLabel("Updating..."),
		//  }
		widget.NewLabel("Processing..."),
		progBar,
		actionButton,
	)

	// create and show the custom dialog
	customDialog = dialog.NewCustomWithoutButtons("Processing Dialog", content, w)
	customDialog.Show()
}
