package ui

import (
	_ "embed"
	"log"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

// used to grab the custom SVG for the heart
//
//go:embed assets/heart.svg
var heartSVG []byte

// creates the toolbar with the options that will be displayed to manage the rendered DB
func createMainWindowToolbar(
	toolbarCanvas fyne.Canvas,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	selectedRow binding.Int,
	dbData *MyDataBinding,
	searchSource binding.String,
	a fyne.App,
) (toolbar *fyne.Container) {
	// create the buttons
	sortButton := createSortButton(sortOrder)
	exportButton := createExportButton(toolbarCanvas)
	settingsButton := createSettingsButton()
	addButton := createAddButton(a, sortCategory, sortOrder, searchText, dbData, selectedRow, searchSource, toolbarCanvas)
	removeButton := createRemoveButton(selectedRow, sortCategory, sortOrder, searchText, dbData)
	helpButton := createHelpButton(toolbarCanvas)
	randButton := createRandomButton(selectedRow, dbData)
	faveButton := createFaveButton(selectedRow, sortCategory, sortOrder, searchText, dbData)
	updateButton := createUpdateButton(sortCategory, sortOrder, searchText, selectedRow, dbData)

	// add them to the toolbar in horizontal row
	toolbar = container.New(
		layout.NewHBoxLayout(),
		sortButton,
		layout.NewSpacer(),
		addButton,
		layout.NewSpacer(),
		updateButton,
		layout.NewSpacer(),
		removeButton,
		layout.NewSpacer(),
		randButton,
		layout.NewSpacer(),
		faveButton,
		layout.NewSpacer(),
		exportButton,
		layout.NewSpacer(),
		settingsButton,
		layout.NewSpacer(),
		helpButton,
	)

	// PERF: change size of each button depending on the size of the given window
	// 1. make toolbar use gridwraplayout
	// toolbar = container.New(
	// 	layout.NewGridWrapLayout(fyne.NewSize(200, 50)),
	// 	sortButton,
	// 	addButton,
	// 	updateButton,
	// 	removeButton,
	// 	randButton,
	// 	faveButton,
	// 	exportButton,
	// 	settingsButton,
	// 	helpButton,
	// )
	// 2. remove text next to buttons

	return toolbar
}

// toggles sort Order (ASC->DESC->ASC)
func createSortButton(sortOrder binding.Bool) (sortButton *widget.Button) {
	// create the button with empty label
	sortButton = widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), func() {
		// whatever curr value of sortOrder is, we want opposite when clicked
		val, _ := sortOrder.Get()
		sortOrder.Set(!val)
	})

	// listen for changes, and update text+icon
	sortOrder.AddListener(binding.NewDataListener(func() {
		val, _ := sortOrder.Get()
		if val {
			sortButton.SetText("Sort ASC")
			sortButton.SetIcon(theme.MenuDropUpIcon())
		} else {
			sortButton.SetText("Sort DESC")
			sortButton.SetIcon(theme.MenuDropDownIcon())
		}
	}))
	return sortButton
}

// export data from db
func createExportButton(toolbarCanvas fyne.Canvas) (exportButton *widget.Button) {
	// create a button without a function
	exportButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), nil)

	// create the dropdown menu items for exporting
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Export to CSV", func() {
			dbhandler.Export(1)
		}),
		fyne.NewMenuItem("Export to SQL", func() {
			dbhandler.Export(2)
		}),
		fyne.NewMenuItem("Export to Markdown", func() {
			dbhandler.Export(3)
		}),
	)

	// define the popup
	menuPopup := widget.NewPopUpMenu(menu, toolbarCanvas)

	// when button clicked, toggle menu
	exportButton.OnTapped = func() {
		menuPopup.ShowAtPosition(exportButton.Position().Add(fyne.NewPos(0, exportButton.Size().Height)))
	}

	return exportButton
}

// TODO: options:
//
//	Light/Dark mode
//	increase/decrease text size
//	update all game data
//	delete entire db
//	default location to store db
//	default location to export to
//	default website to search: HTLB vs Completionator vs Both
//			ss, _ := searchSource.Set("All"/"HLTB"/"COMP")
//
// PERF:
//
//	find way to implement menu without opening new window for window tiling managers
//	Theme Selector
func createSettingsButton() (settingsButton *widget.Button) {
	settingsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		log.Println("Open the dropdown menu and show options for user config")
	})

	return settingsButton
}

// get data and add it to the DB
func createAddButton(
	a fyne.App,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
	searchSource binding.String,
	toolbarCanvas fyne.Canvas,
) (addButton *widget.Button) {
	// no function since we are making a dropdown menu
	addButton = widget.NewButtonWithIcon("Add Game Data", theme.ContentAddIcon(), nil)
	ss, _ := searchSource.Get()
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Single Game Search", func() {
			singleGameNameSearchPopup(
				a,
				searchSource,
				sortCategory,
				sortOrder,
				searchText,
				dbData,
				selectedRow,
			)
		}),
		fyne.NewMenuItem("Manual Entry", func() {
			manualEntryPopup(
				a,
				sortCategory,
				sortOrder,
				searchText,
				dbData,
				selectedRow,
			)
		}),
		fyne.NewMenuItem("From CSV", func() {
			dbhandler.Import(1, ss)
			forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		}),
		// INFO: will drop the existing table and replace with imported SQL file
		fyne.NewMenuItem("From SQL", func() {
			dbhandler.Import(2, ss)
			forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		}),
		// INFO: game names must be separated by new lines with 1 game per line
		fyne.NewMenuItem("From TXT", func() {
			dbhandler.Import(3, ss)
			forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		}),
	)

	menuPopup := widget.NewPopUpMenu(menu, toolbarCanvas)

	addButton.OnTapped = func() {
		menuPopup.ShowAtPosition(addButton.Position().Add(fyne.NewPos(0, addButton.Size().Height)))
	}
	return addButton
}

// finds selected row game name, and deletes it from DB
func createRemoveButton(
	selectedRow binding.Int,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
) (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Remove Game Data", theme.ContentRemoveIcon(), func() {
		selrow, _ := selectedRow.Get()
		if selrow >= 0 {
			// get the game name and send query for deletion
			dbdata, _ := dbData.Get()
			dbhandler.DeleteFromDB(dbdata[selrow][0])

			forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		}
	})

	return removeButton
}

// lists help options such as tutorial, manual, support, etc.
func createHelpButton(toolbarCanvas fyne.Canvas) (helpButton *widget.Button) {
	helpButton = widget.NewButtonWithIcon("", theme.QuestionIcon(), nil)
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Show Tutorial", func() { println("Highlights and focuses what each thing does") }),
		fyne.NewMenuItem("Open Manual", func() {
			println("Opens a new window with booklet/document. Explicit with page numbers and how to do stuff")
		}),
		// TODO: maybe want to have version number here?
		// Update button?
		// Figure out what would go into Program Info and see if i can expand this menu to accommodate that data
		fyne.NewMenuItem("Program Info", func() { println("Opens a new window with information about the Program") }),
		fyne.NewMenuItem("Support Me <3", func() { println("Is a link such that, when clicked will take you to Ko-Fi, Paypal, etc.") }),
	)

	menuPopup := widget.NewPopUpMenu(menu, toolbarCanvas)

	helpButton.OnTapped = func() {
		menuPopup.ShowAtPosition(helpButton.Position().Add(fyne.NewPos(0, helpButton.Size().Height)))
	}
	return helpButton
}

// randomly selects a row to highlight
func createRandomButton(selectedRow binding.Int, dbData *MyDataBinding) (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Random Row", theme.SearchReplaceIcon(), func() {
		dbdata, _ := dbData.Get()
		selectedRow.Set(rand.Intn(len(dbdata)))
	})

	return removeButton
}

// toggle favorite for game defined by selectedRow
func createFaveButton(
	selectedRow binding.Int,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
) (faveButton *widget.Button) {
	heartIcon := fyne.NewStaticResource("heart.svg", heartSVG)
	faveButton = widget.NewButtonWithIcon("(Un)Favorite", theme.NewThemedResource(heartIcon), func() {
		selrow, _ := selectedRow.Get()
		if selrow >= 0 {
			// get the game name and send query for toggling favorite
			dbdata, _ := dbData.Get()
			dbhandler.ToggleFavorite(dbdata[selrow][0])

			forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		}
	})

	return faveButton
}

// update the selected game defined by selectedRow
func createUpdateButton(
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	selectedRow binding.Int,
	dbData *MyDataBinding,
) (updateButton *widget.Button) {
	updateButton = widget.NewButtonWithIcon("Update", theme.MediaReplayIcon(), func() {
		selrow, _ := selectedRow.Get()
		if selrow >= 0 {
			log.Println("Updating highlighted entry")

			forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		}
	})

	return updateButton
}

// TODO: func to update and re-render changes to save space on lines cuz it happens so often
func forceRenderDB(
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) {
	// update dbData and selectedRow to render changes
	updateDBData(sortCategory, sortOrder, searchText, dbData)
	selectedRow.Set(-2)
	selectedRow.Set(-1)
}
