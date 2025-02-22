package ui

import (
	_ "embed"
	// "fmt"
	"log"

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
) (toolbar *fyne.Container) {
	// create the buttons
	sortButton := createSortButton(sortOrder)
	exportButton := createExportButton(toolbarCanvas)
	settingsButton := createSettingsButton()
	addButton := createAddButton(sortOrder, toolbarCanvas)
	removeButton := createRemoveButton(selectedRow, sortCategory, sortOrder, searchText, dbData)
	helpButton := createHelpButton(toolbarCanvas)
	randButton := createRandomButton()
	faveButton := createFaveButton()
	updateButton := createUpdateButton()

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

	return toolbar
}

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

func createExportButton(toolbarCanvas fyne.Canvas) (exportButton *widget.Button) {
	// create a button without a function
	exportButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), nil)

	// create the dropdown menu items for exporting
	// TODO: Export the current view, not the default one in the database
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Export to SQL", func() {
			println("Export to SQL")
			dbhandler.Export(1)
		}),
		fyne.NewMenuItem("Export to CSV", func() {
			println("Export to CSV")
			dbhandler.Export(2)
		}),
		fyne.NewMenuItem("Export to Markdown", func() {
			println("Export to Markdown")
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
//	default website to search first: HTLB vs Completionator
//	dont ask for confirmation for:
//		deleting games
//		updating all entries without URL
//	find way to implement menu without opening new window for window tiling managers
//
// PERF: possible future updates?
//
//	Theme Selector
func createSettingsButton() (settingsButton *widget.Button) {
	settingsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		log.Println("Open the dropdown menu and show options for user config")
	})

	return settingsButton
}

// TODO: connect to dbData
// get data and append the new value
func createAddButton(sortOrder binding.Bool, toolbarCanvas fyne.Canvas) (addButton *widget.Button) {
	addButton = widget.NewButtonWithIcon("Add Game Data", theme.ContentAddIcon(), func() {
		log.Println("dropdown menu of diff ways to add data")
	})

	// TODO: Open window for Manual Entry and Single Game Search
	menu := fyne.NewMenu("",
		// TODO: re render the dbrender widget whenever one of these is called
		fyne.NewMenuItem("Single Game Search", func() {
			println("Open New window for single game name entry")
		}),
		// TODO: re render the dbrender widget whenever one of these is called
		fyne.NewMenuItem("Manual Entry", func() {
			println("Open New Window with form for game data entry")
		}),
		// TODO: Fix the below functions so they re render the DB properly
		// Should be connected to dbData
		// INFO: For TXT File, they must be game names separated by new lines with 1 game per line
		fyne.NewMenuItem("From TXT", func() {
			dbhandler.ImportTXT()
		}),
		// INFO: Cannot guarantee that the program will accept SQL/CSV files not created by this program (i.e. by hand)
		fyne.NewMenuItem("From SQL", func() {
			dbhandler.ImportSQL()
		}),
		fyne.NewMenuItem("From CSV", func() {
			dbhandler.ImportCSV()
		}),
	)

	menuPopup := widget.NewPopUpMenu(menu, toolbarCanvas)

	addButton.OnTapped = func() {
		menuPopup.ShowAtPosition(addButton.Position().Add(fyne.NewPos(0, addButton.Size().Height)))
	}
	return addButton
}

// finds selected row game name, and issues query to db
func createRemoveButton(
	selectedRow binding.Int,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
) (removeButton *widget.Button) {
	// TODO: Ask for confirmation in new window
	removeButton = widget.NewButtonWithIcon("Remove Game Data", theme.ContentRemoveIcon(), func() {
		log.Println("remove the highlighted row")
		selrow, _ := selectedRow.Get()
		if selrow >= 0 {
			// get the game name and send query for deletion
			dbdata, _ := dbData.Get()
			dbhandler.DeleteFromDB(dbdata[selrow][0])

			// update dbData and selectedRow to render changes
			updateDBData(sortCategory, sortOrder, searchText, dbData)
			selectedRow.Set(-1)
		}
	})

	return removeButton
}

func createHelpButton(toolbarCanvas fyne.Canvas) (helpButton *widget.Button) {
	helpButton = widget.NewButtonWithIcon("", theme.QuestionIcon(), nil)
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Show Tutorial", func() { println("Highlights and focuses what each thing does") }),
		fyne.NewMenuItem("Open Manual", func() {
			println("Opens a new window with booklet/document. Explicit with page numbers and how to do stuff")
		}),
		// TODO: maybe want to have version number here?
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

// TODO: Connect this to selectedRow and dbData
// change to random value from 1:rows
func createRandomButton() (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Random Row", theme.SearchReplaceIcon(), func() {
		log.Println("highlight a random row in the other pane")
	})

	return removeButton
}

// TODO: connect to selectedRow and dbData
// get row value from selectedRow and select row from dbData
// toggle favorite for that game in particular
func createFaveButton() (faveButton *widget.Button) {
	heartIcon := fyne.NewStaticResource("heart.svg", heartSVG)
	faveButton = widget.NewButtonWithIcon("(Un)Favorite", theme.NewThemedResource(heartIcon), func() {
		log.Println("(Un)Favorite the selected row")
	})

	return faveButton
}

// TODO:connect to selectedRow and dbData
// get row value and select row from dbData
func createUpdateButton() (updateButton *widget.Button) {
	updateButton = widget.NewButtonWithIcon("Update", theme.MediaReplayIcon(), func() {
		log.Println("updated highlighted entry")
	})

	return updateButton
}
