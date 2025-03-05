package ui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

var prevWidth float32

// makes the table and reflects changes based on values of bindings
// TODO: make the favorited rows be a diff color than the others
func createDBRender(
	selectedRow binding.Int,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	selectedTheme binding.String,
	dbData *MyDataBinding,
	w fyne.Window,
) (dbRender *widget.Table) {
	// given the bindings create the table with the new set of data
	sortcat, _ := sortCategory.Get()
	sortord, _ := sortOrder.Get()

	// if db exists then get the data
	// NOTE: notice there is no search text, because we initialize without any search text from the user
	if dbhandler.CheckDBExists() {
		// no initial search query so use ""
		dbData.Set(dbhandler.SortDB(sortcat, sortord, ""))
	} else {
		dbhandler.CreateDB()
	}
	data, _ := dbData.Get()

	// make the table with size of data. default 1
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	// TODO: Binding for themesDir location
	availableThemes, err := loadAllThemes("themes")
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}

	st, _ := selectedTheme.Get()
	currTheme := availableThemes[st]

	// populate table with info
	dbRender = widget.NewTableWithHeaders(
		// table dims
		func() (int, int) { return numRows, 4 },
		// create empty cells with dflt bg color and empty text
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(hexToColor(currTheme.Background))
			label := widget.NewLabel("")
			return container.NewStack(bg, label)
		},
		// populate table with content
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)

			// if there is data in DB then display it o/w display "No Data"
			if len(data) > 1 {
				if id.Col == 0 {
					// if game name is too long, then truncate and append "...". o/w display entire game name
					if len(data[id.Row][0]) < 48 {
						label.SetText(data[id.Row][0])
					} else {
						label.SetText(data[id.Row][0][:45] + "...")
					}
				} else {
					// display the time data
					label.SetText(fmt.Sprintf("%v", data[id.Row][id.Col]))
				}
			} else {
				label.SetText("No Data")
			}

			// no selected row, so apply background and alt background
			if id.Row%2 == 0 {
				bg.FillColor = hexToColor(currTheme.Background)
			} else {
				bg.FillColor = hexToColor(currTheme.AltBackground)
			}
		},
	)

	// highlight the row of the cell clicked if its not a divider -> dividers have negative position values
	dbRender.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Col >= 0 {
			selectedRow.Set(id.Row)
		}
	}

	// set up the headers
	width := w.Content().Size().Width
	dbRender = headerSetup(sortCategory, selectedTheme, dbRender, width)

	// change contents of dbData binding when sort order changes
	sortOrder.AddListener(binding.NewDataListener(func() {
		updateDBData(sortCategory, sortOrder, searchText, dbData)
		width := w.Content().Size().Width
		dbRender = updateTable(sortCategory, selectedRow, dbData, selectedTheme, dbRender, width)
		forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		dbRender.Refresh()
	}))

	// change contents of dbData binding when sort category changes
	sortCategory.AddListener(binding.NewDataListener(func() {
		updateDBData(sortCategory, sortOrder, searchText, dbData)
		width := w.Content().Size().Width
		dbRender = updateTable(sortCategory, selectedRow, dbData, selectedTheme, dbRender, width)
		forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		dbRender.Refresh()
	}))

	// change contents of dbData binding when search text changes
	searchText.AddListener(binding.NewDataListener(func() {
		updateDBData(sortCategory, sortOrder, searchText, dbData)
		width := w.Content().Size().Width
		dbRender = updateTable(sortCategory, selectedRow, dbData, selectedTheme, dbRender, width)
		forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		dbRender.Refresh()
	}))

	// selectedRow changes
	selectedRow.AddListener(binding.NewDataListener(func() {
		selRow, _ := selectedRow.Get()
		width := w.Content().Size().Width
		dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
			// get the label from the stack
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)
			label.SetText("No Data")

			// highlighting
			if id.Row == selRow {
				bg.FillColor = hexToColor(currTheme.HoverColor)
			} else {
				if id.Row%2 == 0 {
					bg.FillColor = hexToColor(currTheme.Background)
				} else {
					bg.FillColor = hexToColor(currTheme.AltBackground)
				}
			}
		}
		dbRender = updateTable(sortCategory, selectedRow, dbData, selectedTheme, dbRender, width)

		// scroll to the new location selected
		var selCell widget.TableCellID
		selCell.Row = selRow
		selCell.Col = 0
		dbRender.ScrollTo(selCell)
		dbRender.ScrollToLeading()

		dbRender.Refresh()
	}))

	// refresh DBRender when the theme changes
	selectedTheme.AddListener(binding.NewDataListener(func() {
		forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
		dbRender.Refresh()
	}))

	// goroutine to adjust col widths every 0.25 s
	go fixTableSize(sortCategory, selectedRow, dbData, selectedTheme, dbRender, w)

	return
}

// sets dbData with given opts
func updateDBData(
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
) {
	sortcat, _ := sortCategory.Get()
	sortord, _ := sortOrder.Get()
	searchtxt, _ := searchText.Get()

	dbData.Set(dbhandler.SortDB(sortcat, sortord, strings.TrimSpace(searchtxt)))
}

// update the contents of the given table
func updateTable(
	sortCategory binding.String,
	selectedRow binding.Int,
	dbData *MyDataBinding,
	selectedTheme binding.String,
	dbRender *widget.Table,
	width float32,
) *widget.Table {
	selRow, _ := selectedRow.Get()

	// check rows for finding dims
	data, _ := dbData.Get()
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	// TODO: Binding for themesDir location
	availableThemes, err := loadAllThemes("themes")
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}

	st, _ := selectedTheme.Get()

	currTheme := availableThemes[st]

	// set dims
	dbRender.Length = func() (int, int) { return numRows, 4 }
	dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		// get the label from the stack
		stack := obj.(*fyne.Container)
		bg := stack.Objects[0].(*canvas.Rectangle)
		label := stack.Objects[1].(*widget.Label)

		// if there is data in DB then display it
		// o/w display "No Data"
		if len(data) != 0 {
			if id.Col == 0 {
				// if game name is too long, then truncate and append "...". o/w display entire game name
				if len(data[id.Row][0]) < 48 {
					label.SetText(data[id.Row][0])
				} else {
					label.SetText(data[id.Row][0][:45] + "...")
				}
			} else {
				// display the time data
				label.SetText(fmt.Sprintf("%v", data[id.Row][id.Col]))
			}
		} else {
			label.SetText("No Data")
		}

		// highlighting
		if id.Row == selRow {
			bg.FillColor = hexToColor(currTheme.HoverColor)
		} else {
			if id.Row%2 == 0 {
				bg.FillColor = hexToColor(currTheme.Background)
			} else {
				bg.FillColor = hexToColor(currTheme.AltBackground)
			}
		}
	}

	// setup headers
	dbRender = headerSetup(sortCategory, selectedTheme, dbRender, width)
	return dbRender
}

func headerSetup(
	sortCategory binding.String,
	selectedTheme binding.String,
	dbTable *widget.Table,
	width float32,
) *widget.Table {
	// TODO: Binding for themesDir location
	availableThemes, err := loadAllThemes("themes")
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}

	st, _ := selectedTheme.Get()

	currTheme := availableThemes[st]

	// name of each column header
	headers := []string{"Game Name", "Main Story", "Main + Sides", "Completionist"}

	// setup for creating the headers
	dbTable.CreateHeader = func() fyne.CanvasObject {
		return container.NewStack(
			canvas.NewRectangle(hexToColor(currTheme.ScrollBarColor)),
			widget.NewLabelWithStyle(
				"------",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true},
			),
			widget.NewButton("------", nil),
		)
	}

	// make headers display content
	dbTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		containerObj := obj.(*fyne.Container)
		labelBG := containerObj.Objects[0].(*canvas.Rectangle)
		label := containerObj.Objects[1].(*widget.Label)
		button := containerObj.Objects[2].(*widget.Button)

		if id.Row == -1 {
			// display column header buttons
			button.Show()
			labelBG.Hide()
			label.Hide()
			button.SetText(headers[id.Col])
			button.OnTapped = func() {
				// sortCategory gets set to whichever header was clicked
				if id.Col == 0 {
					sortCategory.Set("name")
				} else if id.Col == 1 {
					sortCategory.Set("main")
				} else if id.Col == 2 {
					sortCategory.Set("mainPlus")
				} else {
					sortCategory.Set("comp")
				}
			}
		} else {
			// display row label index, from 1:rows
			button.Hide()
			label.Show()
			labelBG.Show()
			labelBG.FillColor = hexToColor(currTheme.InputBackgroundColor)
			label.SetText(fmt.Sprintf("%d", id.Row+1))
		}
		labelBG.Refresh()
	}

	// set column widths
	dbTable.SetColumnWidth(0, 400)

	// game name has 400, and the row headers take ~70 spacing
	// all other space is to be given to the other columns
	spacing := (width - 400 - 70) / 3

	dbTable.SetColumnWidth(1, spacing)
	dbTable.SetColumnWidth(2, spacing)
	dbTable.SetColumnWidth(3, spacing)
	return dbTable
}

// check window size every 0.25 and adjust size of table if it changes
func fixTableSize(
	sortCategory binding.String,
	selectedRow binding.Int,
	dbData *MyDataBinding,
	selectedTheme binding.String,
	dbRender *widget.Table,
	w fyne.Window,
) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		width := w.Content().Size().Width
		if prevWidth != width {
			select {
			case <-ticker.C:
				prevWidth = width
				updateTable(sortCategory, selectedRow, dbData, selectedTheme, dbRender, width)
				dbRender.Refresh()
			}
		}
	}
}
