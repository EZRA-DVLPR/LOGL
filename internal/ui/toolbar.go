package ui

import (
	_ "embed"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/model"
)

// used to grab the custom SVG for the heart
//
//go:embed assets/heart.svg
var heartSVG []byte

// creates the toolbar with the options that will be displayed to manage the rendered DB
func createMainWindowToolbar(
	availableThemes map[string]ColorTheme,
) (toolbar *fyne.Container) {
	return container.New(
		layout.NewHBoxLayout(),
		createSortButton(),
		layout.NewSpacer(),
		createAddButton(),
		layout.NewSpacer(),
		createUpdateButton(),
		layout.NewSpacer(),
		createRemoveButton(),
		layout.NewSpacer(),
		createRandomButton(),
		layout.NewSpacer(),
		createFaveButton(),
		layout.NewSpacer(),
		createExportButton(),
		layout.NewSpacer(),
		createHelpButton(),
		layout.NewSpacer(),
		createSettingsButton(availableThemes),
		// HACK: just keep this for when I need to do some quick testing
		// layout.NewSpacer(),
		// createTestButton(availableThemes),
	)

	// PERF: remove text next to buttons and leave as option in settings
}

// toggles sort Order (ASC->DESC->ASC)
func createSortButton() (sortButton *widget.Button) {
	// create the button with empty label
	sortButton = widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), func() {
		// whatever curr value of sortOrder is, we want opposite when clicked
		val, _ := model.GetSortOrder()
		log.Println("Sort Order changed to:", !val)
		model.SetSortOrder(!val)
	})

	// listen for changes, and update text+icon
	model.AddSortOrderListener(func(val bool) {
		if val {
			sortButton.SetText("Sort ASC")
			sortButton.SetIcon(theme.MenuDropUpIcon())
		} else {
			sortButton.SetText("Sort DESC")
			sortButton.SetIcon(theme.MenuDropDownIcon())
		}
	})
	return sortButton
}

// export data from db
func createExportButton() (exportButton *widget.Button) {
	menuItems := []*fyne.MenuItem{
		fyne.NewMenuItem("Export to CSV", func() {
			// idea is to have user pick folder, then entry for filename
			dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
				if err != nil {
					log.Println("Error writing to CSV file:", err)
					return
				}
				if uri == nil {
					log.Println("No file Selected to export to CSV")
					return
				}
				defer uri.Close() // close uri when dialog closes

				// delete the partially created file as i will create it in export function
				os.Remove(uri.URI().Path())

				//export
				dbhandler.Export(1, uri.URI().Path())
			}, w)
		}),
		fyne.NewMenuItem("Export to SQL", func() {
			dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
				if err != nil {
					log.Println("Error writing to SQL file:", err)
					return
				}
				if uri == nil {
					log.Println("No file Selected to export to SQL")
					return
				}
				defer uri.Close() // close uri when dialog closes
				os.Remove(uri.URI().Path())
				dbhandler.Export(2, uri.URI().Path())
			}, w)
		}),
		fyne.NewMenuItem("Export to MD", func() {
			dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
				if err != nil {
					log.Println("Error writing to MD file:", err)
					return
				}
				if uri == nil {
					log.Println("No file Selected to export to MD")
					return
				}
				defer uri.Close() // close uri when dialog closes
				os.Remove(uri.URI().Path())
				dbhandler.Export(3, uri.URI().Path())
			}, w)
		}),
	}

	// define the popup
	menuPopup := NewThemeAwareMenu(menuItems, w.Canvas())

	exportButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		menuPopup.Show(exportButton.Position().Add(fyne.NewPos(0, exportButton.Size().Height)))
	})

	// refresh menupopup when the theme changes
	model.AddSelectedThemeListener(func(string) {
		menuPopup.Refresh()
	})

	return exportButton
}

// PERF:
//
//	find way to implement menu without opening new window for window tiling managers
func createSettingsButton(availableThemes map[string]ColorTheme) (settingsButton *widget.Button) {
	settingsButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		settingsPopup(availableThemes)
	})

	return settingsButton
}

// get data and add it to the DB
func createAddButton() (addButton *widget.Button) {
	menuItems := []*fyne.MenuItem{
		fyne.NewMenuItem("Game Search", func() {
			singleGameNameSearchPopup()
		}),
		fyne.NewMenuItem("Manual Entry", func() {
			manualEntryPopup()
		}),
		fyne.NewMenuItem("From CSV", func() {
			fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
				if err != nil {
					log.Println("Error opening CSV file:", err)
					return
				}
				if uri == nil {
					log.Println("No file Selected for importing from CSV")
					return
				}
				defer uri.Close() // close uri when dialog closes
				PopProgressBar(0)
				dbhandler.Import(1, uri.URI().Path())
				UpdateDBData()
			}, w)
			// set file extension to only allow csv files
			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			fileDialog.Show()
		}),
		// INFO: will drop the existing table and replace with imported SQL file
		fyne.NewMenuItem("From SQL", func() {
			fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
				if err != nil {
					log.Println("Error opening SQL file:", err)
					return
				}
				if uri == nil {
					log.Println("No file Selected for importing from SQL")
					return
				}
				defer uri.Close()
				PopProgressBar(2)
				dbhandler.Import(2, uri.URI().Path())
				UpdateDBData()
			}, w)
			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".sql"}))
			fileDialog.Show()
		}),
		// INFO: game names must be separated by new lines with 1 game per line
		fyne.NewMenuItem("From TXT", func() {
			fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
				if err != nil {
					log.Println("Error opening TXT file:", err)
					return
				}
				if uri == nil {
					log.Println("No file Selected for importing from TXT")
					return
				}
				defer uri.Close()
				PopProgressBar(0)
				dbhandler.Import(3, uri.URI().Path())
				UpdateDBData()
			}, w)
			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
			fileDialog.Show()
		}),
		// INFO: These require user input to be completed
		fyne.NewMenuItem("From Epic", func() {
			integrationImport("epic")
		}),
		fyne.NewMenuItem("From GOG", func() {
			integrationImport("gog")
		}),
		fyne.NewMenuItem("From PSN", func() {
			integrationImport("psn")
		}),
		fyne.NewMenuItem("From Steam", func() {
			integrationImport("steam")
		}),
	}

	// define the popup
	menuPopup := NewThemeAwareMenu(menuItems, w.Canvas())

	addButton = widget.NewButtonWithIcon("Add Game", theme.ContentAddIcon(), func() {
		menuPopup.Show(addButton.Position().Add(fyne.NewPos(0, addButton.Size().Height)))
	})

	// refresh menupopup when the theme changes
	model.AddSelectedThemeListener(func(string) {
		menuPopup.Refresh()
	})

	return addButton
}

// finds selected row game name, and deletes it from DB
func createRemoveButton() (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Remove Game", theme.ContentRemoveIcon(), func() {
		selrow, _ := model.GetSelectedRow()
		if selrow >= 0 {
			// get the game name and send query for deletion
			dbdata, _ := dbData.Get()
			log.Println("Removing Game:", dbdata[selrow][0])
			dbhandler.DeleteFromDB(dbdata[selrow][0])

			UpdateDBData()
		}
	})

	return removeButton
}

// lists help options such as tutorial, manual, support, etc.
func createHelpButton() (helpButton *widget.Button) {
	menuItems := []*fyne.MenuItem{
		fyne.NewMenuItem("Video Tutorials", func() {
			goToWebsite("https://youtube.com/playlist?list=PL_gNvZlhoitBNANmcZFgoQpT1FjZiBs7I&si=GBWYIGHiUd0dP2-L")
		}),
		fyne.NewMenuItem("PDF Manual", func() {
			goToWebsite("https://github.com/EZRA-DVLPR/GameList/blob/main/docs/PDF/Manual.pdf")
		}),
		fyne.NewMenuItem("Bug/Feature Tracker", func() {
			goToWebsite("https://github.com/EZRA-DVLPR/GameList/issues")
		}),
		fyne.NewMenuItem("Blog Post", func() {
			goToWebsite("https://personal-website-ezra-dvlpr.vercel.app/blog/projects/GameList")
		}),
		fyne.NewMenuItem("Support Me <3", func() {
			goToWebsite("https://personal-website-ezra-dvlpr.vercel.app/tips")
		}),
	}

	menuPopup := NewThemeAwareMenu(menuItems, w.Canvas())

	helpButton = widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		menuPopup.Show(helpButton.Position().Add(fyne.NewPos(0, helpButton.Size().Height)))
	})

	model.AddSelectedThemeListener(func(string) {
		menuPopup.Refresh()
	})

	return helpButton
}

// randomly selects a row to highlight
func createRandomButton() (removeButton *widget.Button) {
	removeButton = widget.NewButtonWithIcon("Random Row", theme.SearchReplaceIcon(), func() {
		dbdata, _ := dbData.Get()
		model.SetSelectedRow(rand.Intn(len(dbdata)))
	})

	return removeButton
}

// toggle favorite for game defined by selectedRow
func createFaveButton() (faveButton *widget.Button) {
	heartIcon := fyne.NewStaticResource("heart.svg", heartSVG)
	faveButton = widget.NewButtonWithIcon("(Un)Favorite", theme.NewThemedResource(heartIcon), func() {
		selrow, _ := model.GetSelectedRow()
		if selrow >= 0 {
			// get the game name and send query for toggling favorite
			dbdata, _ := dbData.Get()
			dbhandler.ToggleFavorite(dbdata[selrow][0])

			UpdateDBData()
		}
	})

	return faveButton
}

// update the selected game defined by selectedRow
func createUpdateButton() (updateButton *widget.Button) {
	updateButton = widget.NewButtonWithIcon("Update", theme.MediaReplayIcon(), func() {
		selrow, _ := model.GetSelectedRow()
		if selrow >= 0 {
			log.Println("Updating highlighted entry")

			// bring up progress menu
			model.SetMaxProcesses(1)
			PopProgressBar(1)

			dbdata, _ := dbData.Get()
			dbhandler.UpdateGame(dbdata[selrow][0])

			UpdateDBData()
		}
	})

	return updateButton
}

// HACK: just keep this for when I need to do some quick testing
// func createTestButton(availableThemes map[string]ColorTheme) (TestButton *widget.Button) {
// 	TestButton = widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
// 		// anything for testing goes here
// 	})
//
// 	return TestButton
// }

func goToWebsite(link string) {
	var cmd *exec.Cmd

	// change cmd based on which OS is being used
	switch runtime.GOOS {
	case "darwin": // mac = darwin
		cmd = exec.Command("open", link)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", link)
	default: // linux, bsd, etc.
		cmd = exec.Command("xdg-open", link)
	}

	err := cmd.Start()
	if err != nil {
		// Handle error
		log.Println("Error opening link:", link)
		log.Fatal(err)
	}
}
