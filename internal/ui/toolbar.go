package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// creates the toolbar with the options that will be displayed to manage the rendered DB
func createMainWindowToolbar(showLabels bool) (toolbar *fyne.Container) {
	// create the buttons
	sortButton := createSortButton(showLabels)
	exportButton := createExportButton()
	settingsButton := createSettingsButton()
	addButton := createAddButton(showLabels)
	removeButton := createRemoveButton(showLabels)
	helpButton := createHelpButton()
	randButton := createRandomButton(showLabels)
	textSizeUpButton := createTextSizeUpButton(showLabels)
	textSizeDownButton := createTextSizeDownButton(showLabels)
	updateButton := createUpdateButton(showLabels)

	// add them to the toolbar in horizontal row
	toolbar = container.New(
		layout.NewHBoxLayout(),
		sortButton,
		addButton,
		updateButton,
		removeButton,
		randButton,
		textSizeUpButton,
		textSizeDownButton,
		exportButton,
		settingsButton,
		helpButton,
	)

	return toolbar
}

func createSortButton(showLabel bool) (sortButton *widget.Button) {
	isAsc := true

	startText := ""
	if showLabel {
		startText = "Sort ASC"
	}

	sortButton = widget.NewButtonWithIcon(startText, theme.MenuDropUpIcon(), func() {
		if isAsc {
			if showLabel {
				sortButton.SetText("Sort DESC")
			}
			sortButton.SetIcon(theme.MenuDropDownIcon())
			isAsc = false
		} else {
			if showLabel {
				sortButton.SetText("Sort ASC")
			}
			sortButton.SetIcon(theme.MenuDropUpIcon())
			isAsc = true
		}
	})
	return sortButton
}

func createExportButton() (exportButton *widget.Button) {
	exportButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		log.Println("Open the dropdown menu and show options for exporting")
	})

	return exportButton
}

func createSettingsButton() (settingsButton *widget.Button) {
	settingsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		log.Println("Open the dropdown menu and show options for user config")
	})

	return settingsButton
}

func createAddButton(showLabel bool) (addButton *widget.Button) {
	startText := ""
	if showLabel {
		startText = "Add Game Data"
	}

	addButton = widget.NewButtonWithIcon(startText, theme.ContentAddIcon(), func() {
		log.Println("dropdown menu of diff ways to add data")
	})

	return addButton
}

func createRemoveButton(showLabel bool) (removeButton *widget.Button) {
	startText := ""
	if showLabel {
		startText = "Remove Game Data"
	}

	removeButton = widget.NewButtonWithIcon(startText, theme.ContentRemoveIcon(), func() {
		log.Println("remove the highlighted rows")
	})

	return removeButton
}

func createHelpButton() (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		log.Println("opens dropdown of diff options for help menu")
	})

	return removeButton
}

func createRandomButton(showLabel bool) (removeButton *widget.Button) {
	startText := ""
	if showLabel {
		startText = "Random"
	}

	removeButton = widget.NewButtonWithIcon(startText, theme.SearchReplaceIcon(), func() {
		log.Println("highlight a random row in the other pane")
	})

	return removeButton
}

func createTextSizeUpButton(showLabel bool) (textSizeUpButton *widget.Button) {
	startText := ""
	if showLabel {
		startText = "Increase Text Size"
	}

	textSizeUpButton = widget.NewButtonWithIcon(startText, theme.Icon("viewZoomIn"), func() {
		log.Println("adjust text size of everything in app")
	})

	return textSizeUpButton
}

func createTextSizeDownButton(showLabel bool) (textSizeDownButton *widget.Button) {
	startText := ""
	if showLabel {
		startText = "Decrease Text Size"
	}

	textSizeDownButton = widget.NewButtonWithIcon(startText, theme.Icon("viewZoomOut"), func() {
		log.Println("adjust text size of everything in app")
	})

	return textSizeDownButton
}

func createUpdateButton(showLabel bool) (updateButton *widget.Button) {
	startText := ""
	if showLabel {
		startText = "update the entries highlighted"
	}

	updateButton = widget.NewButtonWithIcon(startText, theme.MediaReplayIcon(), func() {
		log.Println("updated highlighted entries")
	})

	return updateButton
}
