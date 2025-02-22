package ui

import (
	"fmt"
	"image/color"
	// "log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

// makes the table and reflects changes based on values of bindings
// TODO: make the favorited rows be a diff color than the others
// PERF: make my own widget (EZRATableWidget) that has the following features:
//  1. get column widths for each column
//  2. set size of column based on size of window
func createDBRender(
	selectedRow binding.Int,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
) (dbRender *widget.Table) {
	// given the bindings create the table with the new set of data
	sortcat, _ := sortCategory.Get()
	sortord, _ := sortOrder.Get()

	// if db exists then get the data
	// NOTE: notice there is no search text, because we initialize without any search text from the user
	if dbhandler.CheckDBExists() {
		// no initial search query so use ""
		dbData.Set(dbhandler.SortDB(sortcat, sortord, ""))
	}
	data, _ := dbData.Get()

	// make the table with size of data. default 1
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	// populate table with info
	dbRender = widget.NewTableWithHeaders(
		// table dims
		func() (int, int) { return numRows, 4 },
		// create empty cells with dflt bg color and empty text
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(color.Black)
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

			// no selected row, so make all cells have bg color
			bg.FillColor = color.Black
		},
	)

	// highlight the row of the cell clicked if its not a divider -> dividers have negative position values
	dbRender.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Col >= 0 {
			selectedRow.Set(id.Row)
		}
	}

	// set up the headers
	dbRender = headerSetup(sortCategory, dbRender)

	// change contents of dbData binding when sort order changes
	sortOrder.AddListener(binding.NewDataListener(func() {
		updateDBData(sortCategory, sortOrder, searchText, dbData)
		dbRender = updateTable(sortCategory, selectedRow, dbData, dbRender)
		selectedRow.Set(-1)
		dbRender.Refresh()
	}))

	// change contents of dbData binding when sort category changes
	sortCategory.AddListener(binding.NewDataListener(func() {
		updateDBData(sortCategory, sortOrder, searchText, dbData)
		dbRender = updateTable(sortCategory, selectedRow, dbData, dbRender)
		selectedRow.Set(-1)
		dbRender.Refresh()
	}))

	// change contents of dbData binding when search text changes
	searchText.AddListener(binding.NewDataListener(func() {
		updateDBData(sortCategory, sortOrder, searchText, dbData)
		dbRender = updateTable(sortCategory, selectedRow, dbData, dbRender)
		selectedRow.Set(-1)
		dbRender.Refresh()
	}))

	// selectedRow changes
	selectedRow.AddListener(binding.NewDataListener(func() {
		selRow, _ := selectedRow.Get()
		dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
			// get the label from the stack
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)
			label.SetText("No Data")

			// highlighting
			if id.Row == selRow {
				bg.FillColor = color.RGBA{0, 0, 180, 255}
			} else {
				bg.FillColor = color.Black
			}
		}
		dbRender = updateTable(sortCategory, selectedRow, dbData, dbRender)

		// scroll to the new location selected
		var selCell widget.TableCellID
		selCell.Row = selRow
		selCell.Col = 0
		dbRender.ScrollTo(selCell)
		dbRender.ScrollToLeading()

		dbRender.Refresh()
	}))

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
	dbRender *widget.Table,
) *widget.Table {
	selRow, _ := selectedRow.Get()

	// check rows for finding dims
	data, _ := dbData.Get()
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

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
			bg.FillColor = color.RGBA{200, 200, 255, 255}
		} else {
			bg.FillColor = color.Black
		}
	}

	// setup headers
	dbRender = headerSetup(sortCategory, dbRender)
	return dbRender
}

func headerSetup(sortCategory binding.String, dbTable *widget.Table) *widget.Table {
	// name of each column header
	headers := []string{"Game Name", "Main Story", "Main + Sides", "Completionist"}

	// setup for creating the headers
	dbTable.CreateHeader = func() fyne.CanvasObject {
		return container.NewStack(
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
		label := containerObj.Objects[0].(*widget.Label)
		button := containerObj.Objects[1].(*widget.Button)

		if id.Row == -1 {
			// display column header buttons
			button.Show()
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
			label.SetText(fmt.Sprintf("%d", id.Row+1))
		}
	}

	// set column widths
	dbTable.SetColumnWidth(0, 400)
	dbTable.SetColumnWidth(1, 200)
	dbTable.SetColumnWidth(2, 200)
	dbTable.SetColumnWidth(3, 200)
	return dbTable
}
