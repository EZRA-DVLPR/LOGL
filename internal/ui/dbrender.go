package ui

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
)

// makes the table and reflects changes based on bindings
// PERF: make my own widget (EZRATableWidget) that has the following features:
//  1. get column widths for each column
//  2. set size of column based on size of window
func createDBRender(sortType binding.String, opt binding.Bool, userText binding.String) (dbRender *widget.Table) {
	// placeholder for data from SQLite requests
	var data [][]string

	// create a new binding that will keep track of the currently selected row
	// INFO: does not persist throughout sessions
	selectedRow := binding.NewInt()
	selectedRow.Set(-1)
	selRow, _ := selectedRow.Get()

	// given the bindings create the table with the new set of data
	sortingType, _ := sortType.Get()
	sortingOpt, _ := opt.Get()
	if sortingOpt {
		data = dbhandler.SortDB(sortingType, "ASC")
	} else {
		data = dbhandler.SortDB(sortingType, "DESC")
	}

	// make the table with size of data. default 1 -- Display "No data"
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	// populate table with info
	dbRender = widget.NewTableWithHeaders(
		// table dims
		func() (int, int) { return numRows, 4 },
		// create an empty cell with def bg color and empty text
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(color.Black)
			label := widget.NewLabel("")
			return container.NewStack(bg, label)
		},
		// populate table with content
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			// get the label from the stack
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			label := stack.Objects[1].(*widget.Label)

			// if there is data in DB then display it
			// o/w display "No Data"
			if numRows > 1 {
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

			// highlighting when the row is selected
			if id.Row == selRow {
				bg.FillColor = color.RGBA{200, 200, 255, 255}
			} else {
				bg.FillColor = color.Black
			}
		},
	)

	// highlight the row if its not a divider -> dividers have negative position values
	dbRender.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Col >= 0 {
			selectedRow.Set(id.Row)
		}
	}

	// set up the header
	dbRender = headerSetup(dbRender)

	// listener to update cell row selection
	selectedRow.AddListener(binding.NewDataListener(func() {
		dbRender = updateTable(selectedRow, opt, sortType, data, dbRender)
		dbRender.Refresh()
	}))

	// listener to update the contents of the table when value of sorting op changes
	opt.AddListener(binding.NewDataListener(func() {
		dbRender = updateTable(selectedRow, opt, sortType, data, dbRender)
		dbRender.Refresh()
	}))

	// listener to update the contents of the table when value of sorting sortType changes
	sortType.AddListener(binding.NewDataListener(func() {
		dbRender = updateTable(selectedRow, opt, sortType, data, dbRender)
		dbRender.Refresh()
	}))

	// listener to update the contents of the table when value of sorting sortType changes
	userText.AddListener(binding.NewDataListener(func() {
		st, _ := userText.Get()
		st = strings.TrimSpace(st)

		if st == "" {
			// return entire DB
		} else {
			// return DB based on game name search with given text
		}
		log.Println("st", st)
	}))
	return
}

// given bindings, data, and table will update the contents of the given table
func updateTable(selectedRow binding.Int, opt binding.Bool, sortBy binding.String, data [][]string, dbRender *widget.Table) *widget.Table {
	sortingType, _ := sortBy.Get()
	sortingOpt, _ := opt.Get()
	selRow, _ := selectedRow.Get()
	if sortingOpt {
		data = dbhandler.SortDB(sortingType, "ASC")
	} else {
		data = dbhandler.SortDB(sortingType, "DESC")
	}
	numRows := 1
	if len(data) != 0 {
		numRows = len(data)
	}

	dbRender.Length = func() (int, int) { return numRows, 4 }
	dbRender.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		// get the label from the stack
		stack := obj.(*fyne.Container)
		bg := stack.Objects[0].(*canvas.Rectangle)
		label := stack.Objects[1].(*widget.Label)

		// if there is data in DB then display it
		// o/w display "No Data"
		if numRows > 1 {
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
	dbRender = headerSetup(dbRender)
	dbRender.ScrollToLeading()
	dbRender.ScrollToTop()
	return dbRender
}

func headerSetup(dbTable *widget.Table) *widget.Table {
	// name of each column header
	headers := []string{"Game Name", "Main Story", "Main + Sides", "Completionist"}

	// setup for creating the headers
	dbTable.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle(
			// add placeholder for at least 6 characters, making 6 digit numbers display nicely
			"______",
			// set text to be centered
			fyne.TextAlignCenter,
			// set text to be bold
			fyne.TextStyle{Bold: true},
		)
	}

	// make headers display content
	dbTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(headers) {
			// display headers defined prev
			obj.(*widget.Label).SetText(headers[id.Col])
		} else {
			// for row index, start at 1:rows
			obj.(*widget.Label).SetText(fmt.Sprintf("%d", id.Row+1))
		}
	}

	// set column widths
	dbTable.SetColumnWidth(0, 400)
	dbTable.SetColumnWidth(1, 200)
	dbTable.SetColumnWidth(2, 200)
	dbTable.SetColumnWidth(3, 200)
	return dbTable
}
