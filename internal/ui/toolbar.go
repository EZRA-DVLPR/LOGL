package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// creates the toolbar with the options that will be displayed to manage the rendered DB
func createMainWindowToolbar(toolbarCanvas fyne.Canvas, sortOrder binding.String) (toolbar *fyne.Container) {
	// create the buttons
	sortButton := createSortButton(sortOrder)
	exportButton := createExportButton(toolbarCanvas)
	settingsButton := createSettingsButton()
	addButton := createAddButton(toolbarCanvas)
	removeButton := createRemoveButton()
	helpButton := createHelpButton(toolbarCanvas)
	randButton := createRandomButton()
	textSizeUpButton := createTextSizeUpButton()
	textSizeDownButton := createTextSizeDownButton()
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
		textSizeUpButton,
		layout.NewSpacer(),
		textSizeDownButton,
		layout.NewSpacer(),
		exportButton,
		layout.NewSpacer(),
		settingsButton,
		layout.NewSpacer(),
		helpButton,
	)

	return toolbar
}

func createSortButton(sortOrder binding.String) (sortButton *widget.Button) {
	// create the button with empty label
	sortButton = widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), func() {
		// whatver curr value of sortOrder is, we want opposite when clicked
		val, _ := sortOrder.Get()
		if val == "ASC" {
			sortOrder.Set("DESC")
		} else {
			sortOrder.Set("ASC")
		}
	})

	// listen for changes, and update text+icon
	sortOrder.AddListener(binding.NewDataListener(func() {
		val, _ := sortOrder.Get()
		if val == "ASC" {
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

	// create the dropdown menu items
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Export to CSV", func() { println("Export to CSV") }),
		fyne.NewMenuItem("Export to SQL", func() { println("Export to SQL") }),
		fyne.NewMenuItem("Export to Markdown", func() { println("Export to Markdown") }),
	)

	// define the popup
	menuPopup := widget.NewPopUpMenu(menu, toolbarCanvas)

	// when button clicked, toggle menu
	exportButton.OnTapped = func() {
		menuPopup.ShowAtPosition(exportButton.Position().Add(fyne.NewPos(0, exportButton.Size().Height)))
	}

	return exportButton
}

// TODO: Might want to open a new window for this one in particular...
func createSettingsButton() (settingsButton *widget.Button) {
	// TODO: special case of adding a stylized menu with diff types of menu options
	// eg. checkbox, slider, filepath, etc.

	// TODO: options:
	//		Light/Dark mode

	// PERF: maybe future updates?
	//		Theme Selector
	settingsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		log.Println("Open the dropdown menu and show options for user config")
	})

	return settingsButton
}

func createAddButton(toolbarCanvas fyne.Canvas) (addButton *widget.Button) {
	addButton = widget.NewButtonWithIcon("Add Game Data", theme.ContentAddIcon(), func() {
		log.Println("dropdown menu of diff ways to add data")
	})

	// TODO: Open window for Manual Entry and Single Game Search
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Single Game Search", func() { println("Open New window for single game name entry") }),
		fyne.NewMenuItem("Manual Entry", func() { println("Open New Window with form for game data entry") }),
		// WARN: For TXT File, they must be game names separated by new lines with 1 game per line
		// Cannot guarantee that the program will accept files not created by this program (i.e. by hand)
		fyne.NewMenuItem("From TXT", func() { println("Import from txt file") }),
		fyne.NewMenuItem("From SQL", func() { println("Import from SQL file") }),
		fyne.NewMenuItem("From CSV", func() { println("Import from CSV file") }),
	)

	menuPopup := widget.NewPopUpMenu(menu, toolbarCanvas)

	addButton.OnTapped = func() {
		menuPopup.ShowAtPosition(addButton.Position().Add(fyne.NewPos(0, addButton.Size().Height)))
	}
	return addButton
}

func createRemoveButton() (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Remove Game Data", theme.ContentRemoveIcon(), func() {
		log.Println("remove the highlighted rows")
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

func createRandomButton() (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Random Row", theme.SearchReplaceIcon(), func() {
		log.Println("highlight a random row in the other pane")
	})

	return removeButton
}

func createTextSizeUpButton() (textSizeUpButton *widget.Button) {
	textSizeUpButton = widget.NewButtonWithIcon("Increase Text Size", theme.Icon("viewZoomIn"), func() {
		log.Println("adjust text size of everything in app")
	})

	return textSizeUpButton
}

func createTextSizeDownButton() (textSizeDownButton *widget.Button) {
	textSizeDownButton = widget.NewButtonWithIcon("Decrease Text Size", theme.Icon("viewZoomOut"), func() {
		log.Println("adjust text size of everything in app")
	})

	return textSizeDownButton
}

func createUpdateButton() (updateButton *widget.Button) {
	updateButton = widget.NewButtonWithIcon("Update", theme.MediaReplayIcon(), func() {
		log.Println("updated highlighted entries")
	})

	return updateButton
}
